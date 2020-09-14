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

type Visitor struct {

	Name string `json:"name" cql:"name"`

	Email string `json:"email,omitempty" cql:"email"`

	Phone string `json:"phone" cql:"phone"`

	// ISO 3166 3-Digit country code
	Country string `json:"country,omitempty" cql:"country"`

	City string `json:"city" cql:"city"`

	ZipCode string `json:"zip_code,omitempty" cql:"zip_code"`

	Street string `json:"street" cql:"street"`
}
