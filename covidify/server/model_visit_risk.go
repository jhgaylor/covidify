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

type VisitRisk struct {

	Risk string `json:"risk,omitempty" cql:"risk"`

	Description string `json:"description,omitempty" cql:"description"`
}
