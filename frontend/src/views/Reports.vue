<script setup lang="ts">
import { ref } from "vue";
import { store } from "../store";
import { GetReport, OpenPath } from "../../wailsjs/go/main/App"; // Adicione OpenPath aqui

const selectedTarget = ref("");
const startDate = ref("");
const startTime = ref("00:00");
const endDate = ref("");
const endTime = ref("23:59");
const isGenerating = ref(false);

// Agora armazenamos os caminhos reais para poder abrir depois
const paths = ref({
  summary: "",
  folder: "",
});

const generateReport = async () => {
  if (!selectedTarget.value || !startDate.value || !endDate.value) {
    alert("Preencha todos os campos.");
    return;
  }

  try {
    isGenerating.value = true;
    paths.value = { summary: "", folder: "" };

    const startFull = `${startDate.value}T${startTime.value}`;
    const endFull = `${endDate.value}T${endTime.value}`;

    // O Go agora deve retornar o caminho do resumo (ajuste seu GetReport para retornar isso)
    const summaryPath = await GetReport(
      selectedTarget.value,
      startFull,
      endFull,
    );

    paths.value.summary = summaryPath;
    // Pega apenas a pasta (removendo o nome do arquivo)
    if (summaryPath) {
      // No Linux/Ubuntu usamos '/'
      paths.value.folder = summaryPath.substring(
        0,
        summaryPath.lastIndexOf("/"),
      );
    }
  } catch (err) {
    alert("Erro: " + err);
  } finally {
    isGenerating.value = false;
  }
};

const openFile = () => {
    if (paths.value.summary) OpenPath(paths.value.summary);
};

const openFolder = () => {
    if (paths.value.folder) OpenPath(paths.value.folder);
};
</script>

<template>
  <div
    class="max-w-4xl mx-auto p-6 bg-gray-900/40 border border-gray-800 rounded-xl"
  >
    <h2 class="text-green-400 font-mono mb-6 uppercase text-sm">
      // SYSTEM EXPORT TERMINAL
    </h2>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="md:col-span-2 space-y-2">
        <label class="text-xs font-mono text-gray-500 uppercase"
          >Target Host</label
        >
        <select
          v-model="selectedTarget"
          class="w-full bg-black border border-gray-700 p-3 rounded-lg text-gray-200 outline-none"
        >
          <option value="" disabled>Select target...</option>
          <option
            v-for="t in store.targets"
            :key="t.id"
            :value="t.id"
            class="text-white"
          >
            {{ t.name }}
          </option>
        </select>
      </div>

      <div class="space-y-4 border-l-2 border-green-500/20 pl-4">
        <label class="text-xs font-mono text-green-500/50 uppercase"
          >Start Point</label
        >
        <input
          type="date"
          v-model="startDate"
          class="w-full bg-black border border-gray-700 p-3 rounded-lg text-gray-200 outline-none block"
        />
        <input
          type="time"
          v-model="startTime"
          class="w-full bg-black border border-gray-700 p-3 rounded-lg text-gray-200 outline-none block"
        />
      </div>

      <div class="space-y-4 border-l-2 border-red-500/20 pl-4">
        <label class="text-xs font-mono text-red-500/50 uppercase"
          >End Point</label
        >
        <input
          type="date"
          v-model="endDate"
          class="w-full bg-black border border-gray-700 p-3 rounded-lg text-gray-200 outline-none block"
        />
        <input
          type="time"
          v-model="endTime"
          class="w-full bg-black border border-gray-700 p-3 rounded-lg text-gray-200 outline-none block"
        />
      </div>

      <button
        @click="generateReport"
        :disabled="isGenerating"
        class="md:col-span-2 w-full bg-green-600 hover:bg-green-500 disabled:bg-gray-800 text-black font-bold py-4 rounded-lg transition-all font-mono text-xs uppercase"
      >
        {{ isGenerating ? "> PROCESSING..." : "> GENERATE REPORTS" }}
      </button>

      <div
        v-if="paths.summary"
        class="md:col-span-2 p-4 bg-green-500/10 border border-green-500/20 rounded-lg space-y-3"
      >
        <div
          class="flex items-center gap-2 text-green-400 font-mono text-[11px]"
        >
          <span class="animate-pulse">●</span> RELATÓRIOS GERADOS COM SUCESSO
        </div>

        <div class="flex flex-wrap gap-3">
          <button
            @click="openFile"
            class="flex-1 bg-green-500/20 hover:bg-green-500/30 text-green-400 border border-green-500/40 py-2 px-4 rounded text-[10px] font-mono transition-all"
          >
            [ ABRIR RESUMO .TXT ]
          </button>

          <button
            @click="openFolder"
            class="flex-1 bg-gray-800 hover:bg-gray-700 text-gray-300 border border-gray-700 py-2 px-4 rounded text-[10px] font-mono transition-all"
          >
            [ ABRIR PASTA DOWNLOADS ]
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
input[type="date"],
input[type="time"],
select {
  color-scheme: dark;
}
</style>
