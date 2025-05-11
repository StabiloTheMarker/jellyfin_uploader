<script setup lang="ts">
import { Button } from "@/components/ui/button";
import { onMounted, ref } from "vue";
import axios, { type AxiosProgressEvent } from "axios";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress"
import { toast, Toaster } from "vue-sonner";
import { mapUploadProcess, type UploadProcess } from "@/types.ts";
import UploadProcessContainer from "@/components/custom/UploadProcessContainer.vue";

const path = ref<string>()
const progress = ref<number | null>(0)
const isUploading = ref<boolean>(false)
const successfullUpload = ref(false)
const files = ref<File[]>([])
const uploadProcesses = ref<UploadProcess[]>([])
const uploadSpeed = ref(0)

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
})
async function handleSubmit() {
  successfullUpload.value = false
  isUploading.value = true
  // First we create the uploadProcess
  const formData = new FormData()
  formData.append("DirPath", path.value as string)
  try {
    const uploadProcessResponse = await axios.post("/api/upload_process", formData, {
      headers: { "Content-Type": "application/json" }
    })
    const data = await uploadProcessResponse.data
    const processId = data.ID
    const fileProgresses: Record<string, number> = {}
    const uploadSpeeds: Record<string, number> = {}
    await Promise.all(files.value.map(it => handleFileUpload(it, processId, fileProgresses, uploadSpeeds)))
    successfullUpload.value = true
  }
  catch (e) {
    toast.error("Could not create Upload Process " + e)
  }
  finally {
    isUploading.value = false
  }
}

async function handleFileUpload(file: File, processId: number, fileProgresses: Record<string, number>, uploadSpeeds: Record<string, number>): Promise<void> {
  const fileName = file.name
  fileProgresses[fileName] = 0
  const fileFormData = new FormData()
  fileFormData.append("files", file)
  let lastTime = Date.now()
  let lastLoaded = 0
  try {
    const response = await axios.post("/api/upload/" + processId, fileFormData, {
      headers: {
        "Content-Type": "multipart/form-data"
      },
      onUploadProgress(progressEvent: AxiosProgressEvent) {

        const now = Date.now()
        const timeElapsed = (now - lastTime) / 1000
        const bytesUploaded = progressEvent.loaded - lastLoaded
        if (timeElapsed > 0) {
          const speedBytesPerSec = bytesUploaded / timeElapsed
          const speedKbps = (speedBytesPerSec / (1024 * 1024))
          uploadSpeeds[fileName] = speedKbps
        }
        uploadSpeed.value = Object.values(uploadSpeeds).reduce((prev, curr) => curr + prev)

        const fileProgress = Math.round((progressEvent.loaded * 100) / (progressEvent.total!!))
        fileProgresses[fileName] = fileProgress
        progress.value = calculateTotalProgress(fileProgresses)
      },
    })
    if (response.status === 200) {
      toast.success("Successfully uploaded file " + fileName)
    }
  }
  catch (e) {
    toast.error("Error in uploading file " + fileName)
  }
}

function calculateTotalProgress(fileProgresses: Record<string, number>): number {
  const count = Object.keys(fileProgresses).length
  const values = Object.values(fileProgresses)
  return values.reduce((prev, curr) => curr + prev) / count
}

</script>

<template>
  <Toaster />
  <div class="max-w-lg m-auto mt-10">
    <h1 class="text-3xl mb-5">Jellyfin Uploader</h1>
    <div class="mb-10 flex flex-col gap-4 max-h-[400px] overflow-y-auto">
      <h2 class="font-semibold text-2xl mb-3">Uploads</h2>
      <UploadProcessContainer @process-deleted="loadProcesses" v-for="process in uploadProcesses"
        :upload-process="process"></UploadProcessContainer>
    </div>
    <div id="container" class="flex flex-col">
      <h2 class="font-semibold text-2xl mb-4">Upload Files</h2>
      <form class="flex flex-col gap-5" @submit.prevent="handleSubmit" enctype="multipart/form-data" method="post">
        <Input v-model="path" class="border border-black" name="path" type="text" />
        <Input @change="(e: Event) => {
          const target = e.target as HTMLInputElement
          files = Array.from(target.files ?? [])
        }" name="files" type="file" multiple />
        <Button type="submit">Submit</Button>
      </form>
      <p v-if="successfullUpload" class="mt-5 text-green-700">Successfully Upload</p>
      <Progress class="mt-5" v-if="isUploading" :model-value="progress"></Progress>
      <div v-if="isUploading" class="flex flex-col mt-5 gap-5">
        <div class="flex justify-between">
          <p>Uploaded Percent: </p>
          <p>{{ progress }}%</p>
        </div>
        <p>Upload Speed: {{ uploadSpeed.toFixed(1) }} MB/s</p>
      </div>
    </div>
  </div>
</template>
