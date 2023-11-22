import { defineStore } from "pinia";

export const usePeopleStore = defineStore({
  id: "people",
  state: () => ({
    showEdit: "",
  }),
  actions: {
    alertType() {},
  },
});
