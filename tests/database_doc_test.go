package tests

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// TestDatabaseDocumentationExists verifies that database.md file exists
func TestDatabaseDocumentationExists(t *testing.T) {
	_, err := os.Stat("../database.md")
	if os.IsNotExist(err) {
		t.Fatal("database.md file does not exist")
	}
	if err != nil {
		t.Fatalf("Error checking database.md: %v", err)
	}
	t.Log("✅ database.md exists")
}

// TestMermaidERDiagramPresent verifies that Mermaid ER Diagram is present
func TestMermaidERDiagramPresent(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)

	// Check for mermaid code block
	if !strings.Contains(contentStr, "```mermaid") {
		t.Fatal("Mermaid code block not found")
	}

	// Check for erDiagram keyword
	if !strings.Contains(contentStr, "erDiagram") {
		t.Fatal("erDiagram keyword not found in Mermaid block")
	}

	t.Log("✅ Mermaid ER Diagram is present")
}

// TestAllEntitiesDefined verifies all 3 entities are defined
func TestAllEntitiesDefined(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)
	entities := []string{"USERS", "TRANSFERS", "POINT_LEDGERS"}

	for _, entity := range entities {
		// Look for entity definition with attributes block
		if !strings.Contains(contentStr, entity+" {") {
			t.Errorf("Entity %s not properly defined with attributes block", entity)
		}
	}

	t.Log("✅ All 3 entities (USERS, TRANSFERS, POINT_LEDGERS) are defined")
}

// TestRelationshipsDefined verifies all relationships are defined
func TestRelationshipsDefined(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)

	// Expected relationships
	relationships := []struct {
		pattern     string
		description string
	}{
		{"USERS ||--o{ TRANSFERS", "USERS to TRANSFERS"},
		{"USERS ||--o{ POINT_LEDGERS", "USERS to POINT_LEDGERS"},
		{"TRANSFERS ||--o{ POINT_LEDGERS", "TRANSFERS to POINT_LEDGERS"},
	}

	for _, rel := range relationships {
		if !strings.Contains(contentStr, rel.pattern) {
			t.Errorf("Relationship '%s' not found", rel.description)
		}
	}

	t.Log("✅ All relationships are properly defined")
}

// TestUserAttributesComplete verifies USERS entity has all required attributes
func TestUserAttributesComplete(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)
	requiredAttributes := []string{
		"id",
		"name",
		"email",
		"phone",
		"address",
		"avatar",
		"points",
		"created_at",
		"updated_at",
		"deleted_at",
	}

	// Extract USERS block
	usersStart := strings.Index(contentStr, "USERS {")
	if usersStart == -1 {
		t.Fatal("USERS entity block not found")
	}

	// Find the closing brace
	usersEnd := strings.Index(contentStr[usersStart:], "}")
	if usersEnd == -1 {
		t.Fatal("USERS entity block not properly closed")
	}

	usersBlock := contentStr[usersStart : usersStart+usersEnd+1]

	for _, attr := range requiredAttributes {
		if !strings.Contains(usersBlock, attr) {
			t.Errorf("USERS entity missing attribute: %s", attr)
		}
	}

	// Check for primary key
	if !strings.Contains(usersBlock, "PK") {
		t.Error("USERS entity missing primary key (PK) designation")
	}

	// Check for unique key on email
	if !strings.Contains(usersBlock, "UK") {
		t.Error("USERS entity missing unique key (UK) designation on email")
	}

	t.Log("✅ USERS entity has all required attributes")
}

// TestTransferAttributesComplete verifies TRANSFERS entity has all required attributes
func TestTransferAttributesComplete(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)
	requiredAttributes := []string{
		"id",
		"from_user_id",
		"to_user_id",
		"amount",
		"status",
		"note",
		"idempotency_key",
		"created_at",
		"updated_at",
		"completed_at",
		"fail_reason",
		"deleted_at",
	}

	transfersStart := strings.Index(contentStr, "TRANSFERS {")
	if transfersStart == -1 {
		t.Fatal("TRANSFERS entity block not found")
	}

	transfersEnd := strings.Index(contentStr[transfersStart:], "}")
	if transfersEnd == -1 {
		t.Fatal("TRANSFERS entity block not properly closed")
	}

	transfersBlock := contentStr[transfersStart : transfersStart+transfersEnd+1]

	for _, attr := range requiredAttributes {
		if !strings.Contains(transfersBlock, attr) {
			t.Errorf("TRANSFERS entity missing attribute: %s", attr)
		}
	}

	// Check for foreign keys
	if strings.Count(transfersBlock, "FK") < 2 {
		t.Error("TRANSFERS entity should have at least 2 foreign keys (from_user_id, to_user_id)")
	}

	t.Log("✅ TRANSFERS entity has all required attributes")
}

// TestPointLedgerAttributesComplete verifies POINT_LEDGERS entity has all required attributes
func TestPointLedgerAttributesComplete(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)
	requiredAttributes := []string{
		"id",
		"user_id",
		"change",
		"balance_after",
		"event_type",
		"transfer_id",
		"reference",
		"metadata",
		"created_at",
	}

	ledgerStart := strings.Index(contentStr, "POINT_LEDGERS {")
	if ledgerStart == -1 {
		t.Fatal("POINT_LEDGERS entity block not found")
	}

	ledgerEnd := strings.Index(contentStr[ledgerStart:], "}")
	if ledgerEnd == -1 {
		t.Fatal("POINT_LEDGERS entity block not properly closed")
	}

	ledgerBlock := contentStr[ledgerStart : ledgerStart+ledgerEnd+1]

	for _, attr := range requiredAttributes {
		if !strings.Contains(ledgerBlock, attr) {
			t.Errorf("POINT_LEDGERS entity missing attribute: %s", attr)
		}
	}

	// Check for foreign keys
	if strings.Count(ledgerBlock, "FK") < 2 {
		t.Error("POINT_LEDGERS entity should have at least 2 foreign keys (user_id, transfer_id)")
	}

	t.Log("✅ POINT_LEDGERS entity has all required attributes")
}

// TestDocumentationSections verifies that all required documentation sections exist
func TestDocumentationSections(t *testing.T) {
	file, err := os.Open("../database.md")
	if err != nil {
		t.Fatalf("Failed to open database.md: %v", err)
	}
	defer file.Close()

	requiredSections := []string{
		"# Database Documentation",
		"## Overview",
		"## Entity Relationship Diagram",
		"## Database Tables",
		"### 1. USERS",
		"### 2. TRANSFERS",
		"### 3. POINT_LEDGERS",
		"## Relationships",
		"## Transaction Flow",
		"## Indexes Strategy",
		"## Data Integrity",
		"## Soft Delete Pattern",
		"## Database Engine",
		"## Migration Strategy",
		"## Query Examples",
		"## Schema Version",
	}

	scanner := bufio.NewScanner(file)
	foundSections := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		for _, section := range requiredSections {
			if strings.HasPrefix(line, section) {
				foundSections[section] = true
			}
		}
	}

	for _, section := range requiredSections {
		if !foundSections[section] {
			t.Errorf("Required section not found: %s", section)
		}
	}

	t.Logf("✅ All %d required documentation sections are present", len(requiredSections))
}

// TestBusinessRulesDocumented verifies that business rules are documented
func TestBusinessRulesDocumented(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)

	businessRules := []string{
		"amount > 0",
		"atomic transaction",
		"idempotency",
		"append-only",
		"soft delete",
		"audit trail",
	}

	for _, rule := range businessRules {
		if !strings.Contains(strings.ToLower(contentStr), strings.ToLower(rule)) {
			t.Errorf("Business rule not documented: %s", rule)
		}
	}

	t.Log("✅ Business rules are documented")
}

// TestMermaidSyntaxValid performs basic Mermaid syntax validation
func TestMermaidSyntaxValid(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)

	// Extract Mermaid block
	mermaidStart := strings.Index(contentStr, "```mermaid")
	if mermaidStart == -1 {
		t.Fatal("Mermaid code block start not found")
	}

	mermaidEnd := strings.Index(contentStr[mermaidStart+10:], "```")
	if mermaidEnd == -1 {
		t.Fatal("Mermaid code block end not found")
	}

	mermaidBlock := contentStr[mermaidStart+10 : mermaidStart+10+mermaidEnd]

	// Check basic syntax
	if !strings.Contains(mermaidBlock, "erDiagram") {
		t.Error("Missing erDiagram declaration")
	}

	// Check for proper relationship syntax (crow's foot notation)
	if !strings.Contains(mermaidBlock, "||--o{") {
		t.Error("Missing proper relationship syntax (||--o{)")
	}

	// Check for entity blocks with braces
	// Note: Relationship syntax uses { as well (e.g., ||--o{), so we only check entity definitions
	entityOpenBraces := 0
	entityCloseBraces := 0
	lines := strings.Split(mermaidBlock, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Count entity definition opening braces (entity name followed by {)
		if strings.Contains(trimmed, "USERS {") || strings.Contains(trimmed, "TRANSFERS {") || strings.Contains(trimmed, "POINT_LEDGERS {") {
			entityOpenBraces++
		}
		// Count standalone closing braces (entity block end)
		if trimmed == "}" {
			entityCloseBraces++
		}
	}
	
	if entityOpenBraces != entityCloseBraces {
		t.Errorf("Unbalanced braces in entity definitions: %d opening, %d closing", entityOpenBraces, entityCloseBraces)
	}

	// Check that all entities have attributes defined
	entities := []string{"USERS", "TRANSFERS", "POINT_LEDGERS"}
	for _, entity := range entities {
		if !strings.Contains(mermaidBlock, entity+" {") {
			t.Errorf("Entity %s missing attribute block definition", entity)
		}
	}

	t.Log("✅ Mermaid syntax validation passed")
}

// TestQueryExamplesProvided verifies that SQL query examples are provided
func TestQueryExamplesProvided(t *testing.T) {
	content, err := os.ReadFile("../database.md")
	if err != nil {
		t.Fatalf("Failed to read database.md: %v", err)
	}

	contentStr := string(content)

	// Check for SQL code blocks
	if !strings.Contains(contentStr, "```sql") {
		t.Error("No SQL query examples found")
	}

	// Check for common query patterns
	queryPatterns := []string{
		"SELECT",
		"WHERE",
		"FROM users",
		"FROM transfers",
		"FROM point_ledgers",
	}

	for _, pattern := range queryPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf("Query pattern not found: %s", pattern)
		}
	}

	t.Log("✅ SQL query examples are provided")
}
