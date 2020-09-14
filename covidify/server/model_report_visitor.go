/*
 * Covidify
 *
 * Simple API collecting guest data.
 *
 * API version: 1.0.0
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package covidify

type ReportVisitor struct {

	Id string `json:"id,omitempty"`

	Visitors []Visitor `json:"visitors,omitempty"`

	Visits []Visit `json:"visits,omitempty"`

	Finalized bool `json:"finalized,omitempty"`

	Contacts []Visitor `json:"contacts,omitempty"`
}
