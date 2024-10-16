<template>
  <v-main>
    <v-tabs>


    </v-tabs>
    <v-form>
      <v-text-field v-model="refs.file.fileName"></v-text-field>
      <v-file-input show-size @change="formChange" v-model="refs.file.fileForm" type='file' name="image"
        label="Image file" variant="solo-inverted" @update="inputFileUpload"></v-file-input>
    </v-form>
    <BtnLoadData @click="inputFileUpload" text="Upload"></BtnLoadData>
  </v-main>
</template>

<script setup>
import { ref } from 'vue'
import BtnLoadData from '@/components/BtnLoadData.vue';

import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";

const config = useConfig;

const refs = ref({
  file: {
    fileName: null,
    fileForm: null,
  }
})

function formChange() {
  console.log("file", refs.value.file.fileForm);
  if (refs.value.file.fileForm) {
    console.log("fileForm empty")
    return
  }

  const regexpSize = /[^0-9a-z\-.]/gi;


  if (refs.value.file.fileForm[0].name) {
    const match = refs.value.file.fileForm[0].name.replace(regexpSize, '');

    refs.value.file.fileName = match
  }

}

async function inputFileUpload() {
  console.log('fileForm', refs.value.file.fileForm)

  if (refs.value.file.fileForm == null) {
    console.log("no files")
    return
  }

  let formData = new FormData();
  formData.append("file", refs.value.file.fileForm[0]);
  // formData.append("file", refs.value.file.fileName);
  let request = {
    method: 'POST',
    body: formData,
  };
  const url = config.BASE_URL + "/tmpl/image";
  await fetchAPI(url, request)
    .then((result) => {
      console.log("file.result", result)
    }).finally(() => {
      console.log("file.result", "finally")
    });

}

</script>
