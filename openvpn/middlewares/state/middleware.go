/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package state

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mysterium/node/openvpn"
	"github.com/mysterium/node/openvpn/management"
)

const stateEventPrefix = ">STATE:"
const stateOutputMatcher = "^\\d+,([a-zA-Z_]+),.*$"

type middleware struct {
	listeners []openvpn.StateCallback
}

// NewMiddleware creates state middleware for given list of callback listeners
func NewMiddleware(listeners ...openvpn.StateCallback) management.Middleware {
	return &middleware{
		listeners: listeners,
	}
}

func (middleware *middleware) Start(commandWriter management.CommandWriter) error {
	middleware.callListeners(openvpn.ProcessStarted)
	_, lines, err := commandWriter.MultiLineCommand("state on all")
	if err != nil {
		return err
	}
	for _, line := range lines {
		state, err := extractOpenvpnState(line)
		if err != nil {
			return err
		}
		middleware.callListeners(state)
	}
	return nil
}

func (middleware *middleware) Stop(commandWriter management.CommandWriter) error {
	middleware.callListeners(openvpn.ProcessExited)
	_, err := commandWriter.SingleLineCommand("state off")
	return err
}

func (middleware *middleware) ConsumeLine(line string) (bool, error) {
	trimmedLine := strings.TrimPrefix(line, stateEventPrefix)
	if trimmedLine == line {
		return false, nil
	}

	state, err := extractOpenvpnState(trimmedLine)
	if err != nil {
		return true, err
	}

	middleware.callListeners(state)
	return true, nil
}

func (middleware *middleware) Subscribe(listener openvpn.StateCallback) {
	middleware.listeners = append(middleware.listeners, listener)
}

func (middleware *middleware) callListeners(state openvpn.State) {
	for _, listener := range middleware.listeners {
		listener(state)
	}
}

func extractOpenvpnState(line string) (openvpn.State, error) {
	rule, err := regexp.Compile(stateOutputMatcher)
	if err != nil {
		return openvpn.UnknownState, err
	}

	matches := rule.FindStringSubmatch(line)
	if len(matches) < 2 {
		return openvpn.UnknownState, errors.New("Line mismatch: " + line)
	}

	return openvpn.State(matches[1]), nil
}
