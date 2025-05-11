<script setup lang="ts">
import { ChevronDown, ChevronUp, Trash } from "lucide-vue-next"
import type { UploadProcess } from "@/types.ts";
import { ref } from "vue";
import { Button } from "@/components/ui/button";

defineProps<{
  uploadProcess: UploadProcess
}>()
const filesOpen = ref(false)

const emits = defineEmits(["processDeleted"])

async function onProcessDeletedClicked(processId: number) {
  const response = await fetch("/api/upload_process/" + processId, {
    method: "DELETE",
  })
  if (!response.ok) {
    throw new Error("Could not delete process with id " + processId)
  }
  else {
    emits("processDeleted")
  }
}
</script>

<template>
  <div class="border border-gray-700 rounded-md p-2">
    <div class="flex justify-between items-center mb-3">
      <p>{{ uploadProcess.CreatedAt.toLocaleString() }}</p>
      <Button @click="() => onProcessDeletedClicked(uploadProcess.ID)" variant="destructive" size="icon"
        class="cursor-pointer">
        <Trash :size="15" />
      </Button>
    </div>
    <div class="flex justify-between items-center">
      <p class="font-semibold">Dateien</p>
      <ChevronDown v-if="!filesOpen" @click="filesOpen = !filesOpen" class="cursor-pointer" />
      <ChevronUp v-else @click="filesOpen = !filesOpen" class="cursor-pointer" />
    </div>
    <div v-if="filesOpen" v-for="file in uploadProcess.Files" class="flex justify-between mt-2 gap-2">
      <p>{{ file.Name }}</p>
      <p class="text-sm">{{ file.Uploaded ? 'Erledigt' : 'Fehlgeschlagen' }}</p>
    </div>
  </div>
</template>
