import { defineStore } from "pinia";

export const userStatuses = ["enabled", "disable", "pending"];

export const userSystemRoles = ["", "admin"];

export class User {
  constructor() {
    this.username = "";
    this.status = "";
    this.tenant_roles;
    this.system_roles;
  }
}

export class EditUser {
  constructor() {
    this.show = false;
    this.user = new User();
    this.isNew = true;
  }
}

export const useUsersStore = defineStore({
  id: "users",
  state: () => ({
    showEdit: "",
  }),
  actions: {
    alertType() {},
  },
});
