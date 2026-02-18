<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { store } from './store'; // Importação da Store
import Dashboard from './views/Dashboard.vue';
import Targets from './views/Targets.vue';

const currentView = ref('dashboard');

onMounted(() => {
    // LISTENER GLOBAL FUNDAMENTAL:
    // Captura os dados do backend e atualiza os números na Store.
    // Sem isso, os cards mostram 0ms mesmo que o gráfico funcione.
    // @ts-ignore
    window.runtime.EventsOn("ping:data", (payload: any) => {
        const latMs = payload.latency / 1000;
        const jitMs = payload.jitter / 1000;
        store.updateStats(payload.hostId, latMs, jitMs, payload.loss);
    });
});
</script>

<template>
  <div class="flex h-screen bg-[#050505] text-gray-300 font-sans selection:bg-green-500/30 selection:text-green-200 overflow-hidden">
    
    <aside class="w-16 md:w-64 bg-black border-r border-gray-900 flex flex-col justify-between shrink-0">
      <div>
        <div class="h-16 flex items-center justify-center border-b border-gray-900">
            <span class="text-green-500 font-bold text-xl tracking-widest">LAG<span class="text-gray-500">MON</span></span>
        </div>
        
        <nav class="mt-8 space-y-2 px-2">
          <button @click="currentView = 'dashboard'" 
            :class="`w-full flex items-center gap-4 px-4 py-3 rounded-lg transition-all ${currentView === 'dashboard' ? 'bg-green-500/10 text-green-400 border border-green-500/20' : 'hover:bg-gray-900 text-gray-500'}`">
            <span>📊</span> <span class="hidden md:block font-mono text-sm">Dashboard</span>
          </button>
          
          <button @click="currentView = 'targets'"
            :class="`w-full flex items-center gap-4 px-4 py-3 rounded-lg transition-all ${currentView === 'targets' ? 'bg-green-500/10 text-green-400 border border-green-500/20' : 'hover:bg-gray-900 text-gray-500'}`">
            <span>🎯</span> <span class="hidden md:block font-mono text-sm">Targets</span>
          </button>
        </nav>
      </div>

      <div class="p-4 text-xs text-center text-gray-700 font-mono">
        v2.1 Stable
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto bg-[#050505] relative">
       <div class="p-6 md:p-8">
          <KeepAlive>
            <component :is="currentView === 'dashboard' ? Dashboard : Targets" />
          </KeepAlive>
       </div>
    </main>
  </div>
</template>