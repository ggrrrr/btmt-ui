import { defineStore } from "pinia";

export class User {
  constructor() {
    this.email = "";
    this.status = "";
    this.SystemRoles = [];
  }
}

export class EditUser {
  constructor() {
    this.show = false;
    this.user = new User();
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
