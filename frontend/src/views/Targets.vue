<script setup lang="ts">
import { ref, onMounted } from "vue";
import { store } from "../store";
import {
  AddTarget,
  RemoveTarget,
  GetTargets,
  SetTargetActive,
  SetTargetDiagramVisibility,
} from "../../wailsjs/go/main/App";

const ip = ref("");
const label = ref("");

// Busca os alvos do backend e popula a store com todos os 6 campos obrigatórios
const loadTargets = async () => {
  const hosts = await GetTargets();
  if (store.targets.length === 0) {
    hosts.forEach((h: any) => 
      store.addTarget(h.id, h.ip, h.name, h.active, h.isGateway, h.showInDiagram)
    );
  }
};

// Adiciona um novo alvo e sincroniza com a store respeitando a nova assinatura
const addNew = async () => {
  if (!ip.value) return;
  const newHost = await AddTarget(ip.value, label.value || ip.value);
  store.addTarget(
    newHost.id, 
    newHost.ip, 
    newHost.name, 
    newHost.active, 
    newHost.isGateway, 
    newHost.showInDiagram
  );
  ip.value = "";
  label.value = "";
};

// Corrige a lógica de ativação: Chama SetTargetActive no backend
const toggleActive = async (target: any) => {
  target.active = !target.active; // Toggle local para feedback imediato
  await SetTargetActive(target.id, target.active);
};

const remove = async (id: string) => {
  await RemoveTarget(id);
  store.removeTarget(id);
};

onMounted(() => {
    loadTargets();
});
</script>

<template>
  <div class="max-w-5xl mx-auto space-y-8">
    <div class="bg-gray-900/50 border border-green-500/20 rounded-xl p-6">
      <h2 class="text-green-400 font-mono text-sm mb-4 uppercase tracking-widest">// Add New Monitor Target</h2>
      <div class="flex gap-4">
        <input v-model="ip" type="text" placeholder="IP Address (e.g. 1.1.1.1)"
          class="flex-1 bg-black border border-gray-700 text-gray-200 p-3 rounded-lg focus:border-green-500 focus:outline-none font-mono text-sm" />

        <input v-model="label" type="text" placeholder="Label (Optional)"
          class="flex-1 bg-black border border-gray-700 text-gray-200 p-3 rounded-lg focus:border-green-500 focus:outline-none font-mono text-sm" />

        <button @click="addNew"
          class="bg-green-600 hover:bg-green-500 text-black font-bold px-6 rounded-lg transition-colors uppercase text-xs tracking-wider">
          Initialize
        </button>
      </div>
    </div>

    <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden">
      <table class="w-full text-left text-sm text-gray-400">
        <thead class="bg-black/50 text-xs uppercase font-mono text-gray-500">
          <tr>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4">Host Details</th>
            <th class="px-6 py-4 text-center">Active Collection</th>
            <th class="px-6 py-4 text-center">Dashboard Graph</th>
            <th class="px-6 py-4 text-center">Diagram Path</th>
            <th class="px-6 py-4 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
          <tr v-for="t in store.targets" :key="t.id" class="hover:bg-gray-800/30 transition-colors">
            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <div :class="`w-2 h-2 rounded-full ${t.active ? (t.stats.loss ? 'bg-red-500 animate-pulse' : 'bg-green-500 shadow-[0_0_8px_#22c55e]') : 'bg-gray-600'}`"></div>
                <span class="font-mono text-[10px]">{{ t.active ? (t.stats.loss ? "LOSS" : "LIVE") : "PAUSED" }}</span>
              </div>
            </td>

            <td class="px-6 py-4">
              <div class="font-bold text-gray-200 flex items-center gap-2">
                {{ t.name }}
                <span v-if="t.isGateway" class="text-[9px] bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 px-1 rounded">GW</span>
              </div>
              <div class="font-mono text-[10px] text-gray-600">{{ t.ip }}</div>
            </td>

            <td class="px-6 py-4 text-center">
              <button @click="toggleActive(t)"
                :class="`w-20 py-1 rounded text-[10px] font-bold border transition-all uppercase tracking-wider ${t.active ? 'bg-green-500/10 text-green-400 border-green-500/30 hover:bg-green-500/20' : 'bg-gray-800 text-gray-500 border-gray-700 hover:border-gray-500'}`">
                {{ t.active ? "ON" : "OFF" }}
              </button>
            </td>

            <td class="px-6 py-4 text-center">
              <label class="inline-flex items-center cursor-pointer relative group">
                <input type="checkbox" v-model="t.showGraph" class="sr-only peer" />
                <div class="w-9 h-5 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-green-600"></div>
              </label>
            </td>

            <td class="px-6 py-4 text-center">
              <input type="checkbox" v-model="t.showInDiagram" @change="SetTargetDiagramVisibility(t.id, t.showInDiagram)"
                class="w-4 h-4 rounded border-gray-700 bg-black text-cyan-600 focus:ring-cyan-500 cursor-pointer" />
            </td>

            <td class="px-6 py-4 text-right">
              <button @click="remove(t.id)" class="text-gray-600 hover:text-red-500 transition-colors p-2 hover:bg-red-500/10 rounded">
                <span class="font-mono text-[10px] font-bold">[DEL]</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>