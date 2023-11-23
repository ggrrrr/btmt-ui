import { defineStore } from "pinia";

export class Dob {
  constructor() {}
}

export class Person {
  constructor() {
    this.id = "";
    this.id_numbers = {};
    this.name = "";
    this.full_name = "";
    this.gender = "";
    this.dob = new Dob();
    this.emails = {};
    this.phones = {};
    this.attr = {};
    this.labels = [];
  }
}

export class EditPerson {
  constructor() {
    this.show = false;
    this.person = new Person();
  }
}

// export default class Person;

export const usePeopleStore = defineStore({
  id: "people",
  state: () => ({
    showEdit: "",
  }),
  actions: {
    alertType() {},
  },
});
