package organization

import (
    "encoding/json"
    "strings"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewOrganization(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        wantErr     bool
        expectedErr string
    }{
        {
            name:    "valid organization name",
            input:   "Valid Organization",
            wantErr: false,
        },
        {
            name:        "empty organization name",
            input:       "",
            wantErr:     true,
            expectedErr: "organization name cannot be empty",
        },
        {
            name:        "whitespace only name",
            input:       "   ",
            wantErr:     true,
            expectedErr: "organization name cannot be empty",
        },
        {
            name:        "extremely long name",
            input:       strings.Repeat("a", 256),
            wantErr:     true,
            expectedErr: "organization name too long",
        },
        {
            name:    "name with special characters",
            input:   "Test & Co. (2023) - Main Branch #1",
            wantErr: false,
        },
        {
            name:    "unicode characters in name",
            input:   "组织名称 🏢 Ørgañızätïøñ",
            wantErr: false,
        },
        {
            name:    "minimum valid length",
            input:   "A",
            wantErr: false,
        },
        {
            name:    "maximum valid length",
            input:   strings.Repeat("a", 255),
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            org, err := NewOrganization(tt.input)

            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedErr)
                assert.Nil(t, org)
            } else {
                require.NoError(t, err)
                require.NotNil(t, org)
                assert.Equal(t, strings.TrimSpace(tt.input), org.Name)
                assert.NotEmpty(t, org.ID)
                assert.False(t, org.CreatedAt.IsZero())
            }
        })
    }
}

func TestOrganization_Validate(t *testing.T) {
    validTime := time.Now()

    tests := []struct {
        name string
        org  Organization
        want bool
    }{
        {
            name: "valid organization",
            org: Organization{
                ID:        "org-123",
                Name:      "Test Org",
                CreatedAt: validTime,
                UpdatedAt: validTime,
            },
            want: true,
        },
        {
            name: "missing ID",
            org: Organization{
                Name:      "Test Org",
                CreatedAt: validTime,
            },
            want: false,
        },
        {
            name: "empty ID",
            org: Organization{
                ID:        "",
                Name:      "Test Org",
                CreatedAt: validTime,
            },
            want: false,
        },
        {
            name: "missing name",
            org: Organization{
                ID:        "org-123",
                CreatedAt: validTime,
            },
            want: false,
        },
        {
            name: "empty name",
            org: Organization{
                ID:        "org-123",
                Name:      "",
                CreatedAt: validTime,
            },
            want: false,
        },
        {
            name: "zero created time",
            org: Organization{
                ID:   "org-123",
                Name: "Test Org",
            },
            want: false,
        },
        {
            name: "future created time",
            org: Organization{
                ID:        "org-123",
                Name:      "Test Org",
                CreatedAt: time.Now().Add(time.Hour),
            },
            want: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.org.Validate()
            assert.Equal(t, tt.want, got)
        })
    }
}

func TestOrganization_String(t *testing.T) {
    org := Organization{
        ID:        "org-123",
        Name:      "Test Organization",
        CreatedAt: time.Now(),
    }

    result := org.String()
    assert.NotEmpty(t, result)
    assert.Contains(t, result, org.Name)
    assert.Contains(t, result, org.ID)
}

func TestOrganization_String_NilReceiver(t *testing.T) {
    var org *Organization
    result := org.String()
    assert.Equal(t, "<nil>", result)
}

func TestOrganization_Equals(t *testing.T) {
    baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

    org1 := Organization{
        ID:        "org-123",
        Name:      "Test Org",
        CreatedAt: baseTime,
    }
    org2 := Organization{
        ID:        "org-123",
        Name:      "Test Org",
        CreatedAt: baseTime,
    }
    org3 := Organization{
        ID:        "org-456",
        Name:      "Test Org",
        CreatedAt: baseTime,
    }

    tests := []struct {
        name string
        org1 Organization
        org2 Organization
        want bool
    }{
        {
            name: "identical organizations",
            org1: org1,
            org2: org2,
            want: true,
        },
        {
            name: "different IDs",
            org1: org1,
            org2: org3,
            want: false,
        },
        {
            name: "same org compared to itself",
            org1: org1,
            org2: org1,
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.org1.Equals(tt.org2)
            assert.Equal(t, tt.want, got)
        })
    }
}

func TestOrganization_JSONMarshaling(t *testing.T) {
    original := Organization{
        ID:        "org-123",
        Name:      "Test Organization",
        CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
        UpdatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
    }

    data, err := json.Marshal(original)
    require.NoError(t, err)
    assert.NotEmpty(t, data)

    var unmarshaled Organization
    err = json.Unmarshal(data, &unmarshaled)
    require.NoError(t, err)
    assert.Equal(t, original.ID, unmarshaled.ID)
    assert.Equal(t, original.Name, unmarshaled.Name)
    assert.True(t, original.CreatedAt.Equal(unmarshaled.CreatedAt))
    assert.True(t, original.UpdatedAt.Equal(unmarshaled.UpdatedAt))
}

func TestOrganization_JSONUnmarshalingInvalidData(t *testing.T) {
    invalidJSONTests := []struct {
        name string
        json string
    }{
        {
            name: "invalid ID type",
            json: `{"id": 123, "name": "Test"}`,
        },
        {
            name: "null name",
            json: `{"id": "123", "name": null}`,
        },
        {
            name: "invalid date format",
            json: `{"id": "123", "name": "Test", "created_at": "invalid-date"}`,
        },
        {
            name: "malformed JSON",
            json: `{"id": "123", "name": "Test"`,
        },
        {
            name: "empty JSON object",
            json: `{}`,
        },
    }
    for _, tt := range invalidJSONTests {
        t.Run(tt.name, func(t *testing.T) {
            var org Organization
            err := json.Unmarshal([]byte(tt.json), &org)
            if err == nil {
                assert.False(t, org.Validate(), "Organization should be invalid after unmarshaling bad JSON")
            }
        })
    }
}

func TestOrganization_Concurrency(t *testing.T) {
    org := &Organization{
        ID:        "org-123",
        Name:      "Test Org",
        CreatedAt: time.Now(),
    }
    const numGoroutines = 10
    const numOperations = 100
    done := make(chan bool, numGoroutines)
    for i := 0; i < numGoroutines; i++ {
        go func() {
            defer func() { done <- true }()
            for j := 0; j < numOperations; j++ {
                _ = org.String()
                _ = org.Validate()
                _ = org.ID
                _ = org.Name
            }
        }()
    }
    for i := 0; i < numGoroutines; i++ {
        select {
        case <-done:
        case <-time.After(5 * time.Second):
            t.Fatal("Timeout waiting for goroutines to complete")
        }
    }
}

func BenchmarkOrganization_Validate(b *testing.B) {
    org := Organization{
        ID:        "org-123",
        Name:      "Test Organization",
        CreatedAt: time.Now(),
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        org.Validate()
    }
}

func BenchmarkOrganization_String(b *testing.B) {
    org := Organization{
        ID:        "org-123",
        Name:      "Test Organization",
        CreatedAt: time.Now(),
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = org.String()
    }
}

func BenchmarkOrganization_JSONMarshal(b *testing.B) {
    org := Organization{
        ID:        "org-123",
        Name:      "Test Organization",
        CreatedAt: time.Now(),
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = json.Marshal(org)
    }
}

func TestOrganization_EdgeCases(t *testing.T) {
    t.Run("nil organization methods", func(t *testing.T) {
        var org *Organization
        assert.NotPanics(t, func() { _ = org.String() })
        assert.NotPanics(t, func() { _ = org.Validate() })
    })
    t.Run("empty struct", func(t *testing.T) {
        org := Organization{}
        assert.False(t, org.Validate())
        assert.NotEmpty(t, org.String())
    })
    t.Run("very long organization name handling", func(t *testing.T) {
        longName := strings.Repeat("A", 10000)
        org := Organization{
            ID:        "org-123",
            Name:      longName,
            CreatedAt: time.Now(),
        }
        assert.NotPanics(t, func() { _ = org.String() })
        assert.NotPanics(t, func() { _ = org.Validate() })
    })
    t.Run("boundary dates", func(t *testing.T) {
        tests := []struct {
            name      string
            createdAt time.Time
            valid     bool
        }{
            {
                name:      "unix epoch",
                createdAt: time.Unix(0, 0),
                valid:     true,
            },
            {
                name:      "far future date",
                createdAt: time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
                valid:     false,
            },
            {
                name:      "current time",
                createdAt: time.Now(),
                valid:     true,
            },
        }
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                org := Organization{
                    ID:        "org-123",
                    Name:      "Test Org",
                    CreatedAt: tt.createdAt,
                }
                assert.Equal(t, tt.valid, org.Validate())
            })
        }
    })
}

func FuzzOrganizationName(f *testing.F) {
    testCases := []string{
        "Test Organization",
        "",
        "A",
        "🏢 Unicode Org",
        "Test & Co.",
        strings.Repeat("a", 100),
    }
    for _, tc := range testCases {
        f.Add(tc)
    }
    f.Fuzz(func(t *testing.T, name string) {
        org, err := NewOrganization(name)
        if err == nil {
            require.NotNil(t, org)
            assert.True(t, org.Validate(), "Created organization should be valid")
            assert.NotEmpty(t, org.ID, "Created organization should have an ID")
            assert.Equal(t, strings.TrimSpace(name), org.Name)
        }
        if org != nil {
            assert.NotPanics(t, func() { _ = org.String() })
            assert.NotPanics(t, func() { _, _ = json.Marshal(org) })
        }
    })
}

func createValidTestOrganization(name string) Organization {
    id := "test-" + strings.ToLower(strings.ReplaceAll(name, " ", "-"))
    return Organization{
        ID:        id,
        Name:      name,
        CreatedAt: time.Now().Add(-time.Hour),
        UpdatedAt: time.Now(),
    }
}

func assertOrganizationValid(t *testing.T, org Organization) {
    t.Helper()
    assert.True(t, org.Validate(), "Organization should be valid")
    assert.NotEmpty(t, org.ID, "Organization should have an ID")
    assert.NotEmpty(t, org.Name, "Organization should have a name")
    assert.False(t, org.CreatedAt.IsZero(), "Organization should have a created time")
}

func assertOrganizationEquals(t *testing.T, got, want Organization) {
    t.Helper()
    assert.Equal(t, want.ID, got.ID, "ID should match")
    assert.Equal(t, want.Name, got.Name, "Name should match")
    assert.True(t, want.CreatedAt.Equal(got.CreatedAt), "CreatedAt should match")
    assert.True(t, want.UpdatedAt.Equal(got.UpdatedAt), "UpdatedAt should match")
}

func TestOrganization_CompleteWorkflow(t *testing.T) {
    tests := []struct {
        name     string
        orgName  string
        wantErr  bool
        testFunc func(t *testing.T, org *Organization)
    }{
        {
            name:    "complete valid workflow",
            orgName: "Test Company Inc.",
            wantErr: false,
            testFunc: func(t *testing.T, org *Organization) {
                assert.True(t, org.Validate())
                str := org.String()
                assert.Contains(t, str, org.Name)
                data, err := json.Marshal(org)
                require.NoError(t, err)
                var unmarshaled Organization
                err = json.Unmarshal(data, &unmarshaled)
                require.NoError(t, err)
                assert.True(t, org.Equals(unmarshaled))
            },
        },
        {
            name:    "workflow with special characters",
            orgName: "Ørg & Co. (2023) #1 🏢",
            wantErr: false,
            testFunc: func(t *testing.T, org *Organization) {
                assert.True(t, org.Validate())
                assert.NotPanics(t, func() { _ = org.String() })
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            org, err := NewOrganization(tt.orgName)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            require.NoError(t, err)
            require.NotNil(t, org)
            if tt.testFunc != nil {
                tt.testFunc(t, org)
            }
        })
    }
}