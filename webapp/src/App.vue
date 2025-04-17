<script setup lang="ts">
import {Button} from "@/components/ui/button";
import {ref} from "vue";
import axios, {type AxiosProgressEvent} from "axios";
import {Input} from "@/components/ui/input";
import {Progress} from "@/components/ui/progress"
const path = ref<string>()
const inputRef = ref<HTMLInputElement | null>(null)
const progress = ref<number | null>(0)
const isUploading = ref<boolean>(true)

async function handleSubmit() {
  const formData = new FormData()
  formData.append("path", path.value as string)
  if (inputRef.value !== null && inputRef.value.files !== null) {
    for (const file of inputRef.value.files) {
      formData.append("files", file)
    }
  }

  try {
    const response = await axios.post("/api/upload", formData, {
      headers: {"Content-Type": "multipart/form-data"},
      onUploadProgress: (progressEvent: AxiosProgressEvent) => {
        if (progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          progress.value = percentCompleted
          console.log(`Upload progress: ${percentCompleted}%`)
        } else {
          console.log(`Uploaded ${progressEvent.loaded} bytes (total unknown)`)
        }
      }
    })
    console.log("response status", response.status)
  } catch (error) {
    console.error("There was an error", error)
  }
}

</script>

<template>
  <div id="container" class="h-screen w-screen flex justify-center items-center">
    <div class="flex flex-col">
      <form class="flex flex-col gap-5" @submit.prevent="handleSubmit" enctype="multipart/form-data" method="post">
        <Input v-model="path" class="border border-black" name="path" type="text"/>
        <Input ref="inputRef" name="files" type="file" multiple/>
        <Button type="submit">Submit</Button>
      </form>
      <Progress class="mt-5" v-if="isUploading" :model-value="progress"></Progress>
    </div>
  </div>
</template>