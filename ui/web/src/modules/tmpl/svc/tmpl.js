import { parseTimestamp } from "@/store/svc";

export class Template {
  constructor(row) {
    if (row === undefined) {
      this.id;
      this.name;
      this.content_type;
      this.body = "";
      this.labels;
      this.editUrl = "";
      this.created_at = null;
      return;
    }
    console.log("", row);
    this.id = row.id;
    this.name = row.name;
    this.content_type = row.content_type;
    this.body = row.body;
    this.labels = row.labels;
    this.editUrl = `/templates/manage/edit/${row.id}`;
    this.created_at = parseTimestamp(row.created_at);
  }
}
