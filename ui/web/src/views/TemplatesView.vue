<template>
  <v-main>
    <v-form>
      <v-text-field v-model="refs.file.fileName"></v-text-field>
      <v-file-input accept="image/*" show-size @change="formChange" v-model="refs.file.fileForm" type='file'
        name="image" label="Image file" variant="solo-inverted" @update="inputFileUpload"></v-file-input>
    </v-form>
    <BtnLoadData @click="inputFileUpload" text="Upload"></BtnLoadData>
  </v-main>
</template>

<script setup>
import { ref } from 'vue'
import BtnLoadData from '@/components/BtnLoadData.vue';
import { useLoginStore } from "@/store/auth";
const loginStore = useLoginStore();


const refs = ref({
  file: {
    fileName: null,
    fileForm: null,
  }
})

function formChange() {
  const regexpSize = /[^0-9a-z\-.]/gi;

  console.log("file", refs.value.file.fileForm[0]);

  if (refs.value.file.fileForm[0].name) {
    const match = refs.value.file.fileForm[0].name.replace(regexpSize, '');

    refs.value.file.fileName = match
  }

}

function inputFileUpload() {
  console.log('inputFileChange', refs.value.file.fileForm[0])
  console.log('inputFileChange', refs.value.file.fileName)

  let headers = {}
  let formData = new FormData();

  if (loginStore.token) {
    headers["Authorization"] = "Bearer " + loginStore.token;
  }

  formData.append("file", refs.value.file.fileForm[0]);
  formData.append("file", refs.value.file.fileName);
  let options = {
    method: 'POST',
    body: formData,
    headers: headers,
  };
  fetch('http://localhost:8010/tmpl/image', options).then(response => response.json())
    .then(data => console.log(data));
}

</script>
