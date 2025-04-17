<script setup lang="ts">
import {Button} from "@/components/ui/button";
import {ref} from "vue";
import axios, {type AxiosProgressEvent} from "axios";
import {Input} from "@/components/ui/input";
import {Progress} from "@/components/ui/progress"
import {toast, Toaster} from "vue-sonner";

const path = ref<string>()
const progress = ref<number | null>(0)
const isUploading = ref<boolean>(false)
const files = ref<File[]>([])

async function handleSubmit() {
  const formData = new FormData()
  formData.append("path", path.value as string)
  for (const file of files.value) {
    formData.append("files", file)
  }

  try {
    isUploading.value = true
    await axios.post("/api/upload", formData, {
      headers: {"Content-Type": "multipart/form-data"},
      onUploadProgress: (progressEvent: AxiosProgressEvent) => {
        if (progressEvent.total) {
          progress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        } else {
          console.log(`Uploaded ${progressEvent.loaded} bytes (total unknown)`)
        }
      }
    })
    toast.success("Upload successfully.")
  } catch (error) {
    console.error("There was an error", error)
    toast.error("There was an error uploading the file")
  } finally {
    isUploading.value = false
  }
}

</script>

<template>
  <Toaster/>
  <div id="container" class="h-screen w-screen flex justify-center items-center">
    <div class="flex flex-col">
      <form class="flex flex-col gap-5" @submit.prevent="handleSubmit" enctype="multipart/form-data" method="post">
        <Input v-model="path" class="border border-black" name="path" type="text"/>
        <Input @change="(e: Event) => {
          const target = e.target as HTMLInputElement
          files = Array.from(target.files ?? [])
        }" name="files" type="file" multiple/>
        <Button type="submit">Submit</Button>
      </form>
      <Progress class="mt-5" v-if="isUploading" :model-value="progress"></Progress>
    </div>
  </div>
</template>