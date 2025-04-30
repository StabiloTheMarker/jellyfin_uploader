<script setup lang="ts">
import {Button} from "@/components/ui/button";
import {onMounted, ref} from "vue";
import axios, {type AxiosProgressEvent} from "axios";
import {Input} from "@/components/ui/input";
import {Progress} from "@/components/ui/progress"
import {toast, Toaster} from "vue-sonner";
import {mapUploadProcess, type UploadProcess} from "@/types.ts";
import UploadProcessContainer from "@/components/custom/UploadProcessContainer.vue";

const path = ref<string>()
const progress = ref<number | null>(0)
const isUploading = ref<boolean>(false)
const successfullUpload = ref(false)
const files = ref<File[]>([])
const uploadProcesses = ref<UploadProcess[]>([])

async function loadUploadProcesses(): Promise<UploadProcess[]> {
  const response = await axios.get("/api/upload_process")
  if (response.status === 200) {
    const data = await response.data
    return data.map(mapUploadProcess);
  } else {
    throw new Error("Could not load Processes")
  }
}

async function loadProcesses() {
  uploadProcesses.value = await loadUploadProcesses()
}

onMounted(async () => {
  await loadProcesses()
  console.log(uploadProcesses.value)
})

async function handleSubmit() {
  successfullUpload.value = false
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
          loadProcesses()
        } else {
          console.log(`Uploaded ${progressEvent.loaded} bytes (total unknown)`)
        }
      }
    })
    toast.success("Upload successfully.")
    successfullUpload.value = true;
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
  <div class="max-w-lg m-auto mt-10">
    <h1 class="text-3xl mb-5">Jellyfin Uploader</h1>
    <div class="mb-10 flex flex-col gap-4 max-h-[400px] overflow-y-auto">
      <h2 class="font-semibold text-2xl mb-3">Uploads</h2>
      <UploadProcessContainer @process-deleted="loadProcesses" v-for="process in uploadProcesses" :upload-process="process"></UploadProcessContainer>
    </div>
    <div id="container" class="flex flex-col">
      <h2 class="font-semibold text-2xl mb-4">Upload Files</h2>
      <form class="flex flex-col gap-5" @submit.prevent="handleSubmit" enctype="multipart/form-data" method="post">
        <Input v-model="path" class="border border-black" name="path" type="text"/>
        <Input @change="(e: Event) => {
          const target = e.target as HTMLInputElement
          files = Array.from(target.files ?? [])
        }" name="files" type="file" multiple/>
        <Button type="submit">Submit</Button>
      </form>
      <p v-if="successfullUpload" class="mt-5 text-green-700">Successfully Upload</p>
      <Progress class="mt-5" v-if="isUploading" :model-value="progress"></Progress>
      <div v-if="isUploading" class="flex mt-5 justify-between">
        <p>Uploaded Percent: </p>
        <p>{{ progress }}%</p>
      </div>
    </div>
  </div>
</template>