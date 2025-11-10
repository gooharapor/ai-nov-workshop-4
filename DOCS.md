## ขอ ER Diagram ใน Mermaid format
- ขอ Database Document follow ตาม practice https://docs.mermaidchart.com/mermaid-oss/syntax/entityRelationshipDiagram.html
- update ผ่าน repo ชื่อ `database.md` ผ่าน git มาได้เลย​โดยใช้ tag mermaid เป็นตัวกำกับ เช่น

```mermaid
erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains
    CUSTOMER }|..|{ DELIVERY-ADDRESS : uses