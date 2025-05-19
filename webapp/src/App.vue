<script setup lang="ts">
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label"
import { computed, onMounted, ref } from "vue";
import axios, { type AxiosProgressEvent } from "axios";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress"
import { toast, Toaster } from "vue-sonner";
import { mapUploadProcess, type UploadProcess } from "@/models.ts";
import UploadProcessContainer from "@/components/custom/UploadProcessContainer.vue";
import { useIntervalFn } from "@vueuse/core";

const path = ref<string>()
const isUploading = ref<boolean>(false)
const successfullUpload = ref(false)
const files = ref<File[]>([])
const uploadProcesses = ref<UploadProcess[]>([])
const uploadedBytes = ref<number>(0)
const totalBytes = ref<number>(0)
const lastUploadedBytes = ref(0)
const uploadSpeedMBps = ref(0)
const etaMinutes = ref(0)
const progressPercent = computed(() => {
  if (totalBytes.value === 0) return 0
  return (uploadedBytes.value / totalBytes.value) * 100
})
useIntervalFn(() => {
  const diff = uploadedBytes.value - lastUploadedBytes.value
  uploadSpeedMBps.value = diff / 1024 / 1024 // convert to MB/s

  const remainingBytes = totalBytes.value - uploadedBytes.value
  const secondsLeft = uploadSpeedMBps.value > 0
    ? remainingBytes / (uploadSpeedMBps.value * 1024 * 1024)
    : 0

  etaMinutes.value = secondsLeft / 60
  lastUploadedBytes.value = uploadedBytes.value
}, 1000)


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

  const formData = new FormData()
  formData.append("DirPath", path.value as string)
  try {
    const uploadProcessResponse = await axios.post("/api/upload_process", formData, {
      headers: { "Content-Type": "application/json" }
    })
    const data = await uploadProcessResponse.data
    const processId = data.ID
    const fileFormData = new FormData()
    for (const file of files.value) {
      fileFormData.append("files[]", file)
    }
    try {
      await axios.post("/api/upload/" + processId, fileFormData, {
        headers: {
          "Content-Type": "multipart/form-data"
        },
        onUploadProgress(progressEvent: AxiosProgressEvent) {
          uploadedBytes.value = progressEvent.loaded
          totalBytes.value = progressEvent.total ?? totalBytes.value
          loadUploadProcesses()
        }
      })
      successfullUpload.value = true
    }
    catch (e) {
      successfullUpload.value = false
      toast.error({ message: e })
    }
  }
  catch (e) {
    toast.error("Could not create Upload Process " + e)
  }
  finally {
    isUploading.value = false
  }
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
        <div>
          <Label class="mb-2">Directory Path</Label>
          <Input v-model="path" class="border border-black" name="path" type="text" />
        </div>
        <div>
          <Label class="mb-2">Files to Upload</Label>
          <Input @change="(e: Event) => {
            const target = e.target as HTMLInputElement
            files = Array.from(target.files ?? [])
          }" name="files" type="file" multiple />
        </div>
        <Button type="submit">Submit</Button>
      </form>
      <p v-if="successfullUpload" class="mt-5 text-green-700">Successfully Upload</p>
      <Progress class="mt-5" v-if="isUploading" :model-value="progressPercent"></Progress>
      <div v-if="isUploading" class="flex flex-col mt-5 gap-5">
        <div class="flex justify-between">
          <p>Uploaded Percent: </p>
          <p>{{ progressPercent?.toFixed(2) }}%</p>
        </div>
        <div class="flex justify-between">
          <p>Upload Speed: {{ uploadSpeedMBps.toFixed(1) }} MB/s</p>
          <p>ETA in Minutes: {{ etaMinutes.toFixed(2) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
