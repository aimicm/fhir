// Copyright (c) 2011-2015, HL7, Inc & The MITRE Corporation
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
//     * Redistributions of source code must retain the above copyright notice, this
//       list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright notice,
//       this list of conditions and the following disclaimer in the documentation
//       and/or other materials provided with the distribution.
//     * Neither the name of HL7 nor the names of its contributors may be used to
//       endorse or promote products derived from this software without specific
//       prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package models

import "encoding/json"

type EligibilityRequest struct {
	Id              string        `json:"-" bson:"_id"`
	Identifier      []Identifier  `bson:"identifier,omitempty" json:"identifier,omitempty"`
	Ruleset         *Coding       `bson:"ruleset,omitempty" json:"ruleset,omitempty"`
	OriginalRuleset *Coding       `bson:"originalRuleset,omitempty" json:"originalRuleset,omitempty"`
	Created         *FHIRDateTime `bson:"created,omitempty" json:"created,omitempty"`
	Target          *Reference    `bson:"target,omitempty" json:"target,omitempty"`
	Provider        *Reference    `bson:"provider,omitempty" json:"provider,omitempty"`
	Organization    *Reference    `bson:"organization,omitempty" json:"organization,omitempty"`
}

type EligibilityRequestBundle struct {
	Id    string                          `json:"id,omitempty"`
	Type  string                          `json:"resourceType,omitempty"`
	Base  string                          `json:"base,omitempty"`
	Total int                             `json:"total,omitempty"`
	Link  []BundleLinkComponent           `json:"link,omitempty"`
	Entry []EligibilityRequestBundleEntry `json:"entry,omitempty"`
}

type EligibilityRequestBundleEntry struct {
	Id       string                `json:"id,omitempty"`
	Base     string                `json:"base,omitempty"`
	Link     []BundleLinkComponent `json:"link,omitempty"`
	Resource EligibilityRequest    `json:"resource,omitempty"`
}

func (resource *EligibilityRequest) MarshalJSON() ([]byte, error) {
	x := struct {
		ResourceType string `json:"resourceType"`
		EligibilityRequest
	}{
		ResourceType:       "EligibilityRequest",
		EligibilityRequest: *resource,
	}
	return json.Marshal(x)
}