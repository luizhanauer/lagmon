<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { GetDiagramConfig } from '../../wailsjs/go/main/App'; // Importe a função

// Estado reativo baseado na configuração
const config = ref({
    local: { name: "Local", ip: "127.0.0.1" },
    gateway: { name: "Gateway", ip: "" },
    internet: { name: "Internet", ip: "" }
});

const nodes = ref({
    local: { status: "ok", lat: 0 },
    gateway: { status: "pending", lat: 0 },
    internet: { status: "pending", lat: 0 }
});

const updateNode = (id: 'gateway' | 'internet', latency: number, loss: boolean) => {
    const node = nodes.value[id];
    node.lat = latency / 1000;
    
    if (loss) {
        node.status = "error";
    } else if (node.lat > 100) {
        node.status = "warn";
    } else {
        node.status = "ok";
    }
};

onMounted(async () => {
    // 1. Busca configuração do JSON
    const data = await GetDiagramConfig();
    config.value = data;

    // 2. Escuta eventos
    // @ts-ignore
    window.runtime.EventsOn("ping:data", (payload: any) => {
        // Compara com os IPs configurados no JSON em vez de IDs hardcoded
        if (payload.ip === config.value.gateway.ip) {
            updateNode("gateway", payload.latency, payload.loss);
        } else if (payload.ip === config.value.internet.ip) {
            updateNode("internet", payload.latency, payload.loss);
        }
    });
});

const getStatusColor = (status: string) => {
    switch(status) {
        case 'error': return 'bg-red-500 shadow-[0_0_15px_rgba(239,68,68,0.6)]';
        case 'warn': return 'bg-yellow-400 shadow-[0_0_15px_rgba(250,204,21,0.6)]';
        default: return 'bg-cyan-500 shadow-[0_0_15px_rgba(6,182,212,0.6)]';
    }
};
</script>

<template>
    <div class="flex items-center justify-between w-full max-w-4xl p-6 bg-black/40 rounded-xl border border-gray-800 backdrop-blur-sm">
        
        <div class="flex flex-col items-center z-10">
            <div class="w-16 h-16 rounded-full flex items-center justify-center bg-gray-700 border-2 border-gray-500">
                <span class="text-2xl">💻</span>
            </div>
            <span class="mt-2 text-xs font-mono text-gray-400 uppercase">{{ config.local.name }}</span>
        </div>

        <div class="flex-1 h-1 bg-gray-800 mx-2 relative overflow-hidden">
            <div class="absolute inset-0 bg-cyan-500/20 animate-pulse"></div>
            <div class="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono text-cyan-400">LAN</div>
        </div>

        <div class="flex flex-col items-center z-10 transition-all duration-300">
            <div :class="`w-16 h-16 rounded-full flex items-center justify-center border-2 border-transparent text-black font-bold transition-all duration-300 ${getStatusColor(nodes.gateway.status)}`">
                <span v-if="nodes.gateway.status === 'pending'">...</span>
                <span v-else>{{ nodes.gateway.lat.toFixed(0) }}ms</span>
            </div>
            <span class="mt-2 text-xs font-mono text-gray-400 uppercase">{{ config.gateway.name }}</span>
            <span class="text-[10px] text-gray-600 font-mono">{{ config.gateway.ip }}</span>
        </div>

        <div class="flex-1 h-1 bg-gray-800 mx-2 relative">
             <div :class="`w-full h-full transition-colors duration-500 ${nodes.gateway.lat > 50 ? 'bg-yellow-500/50' : 'bg-cyan-500/20'}`"></div>
        </div>

        <div class="flex flex-col items-center z-10">
            <div :class="`w-16 h-16 rounded-full flex items-center justify-center border-2 border-transparent text-black font-bold transition-all duration-300 ${getStatusColor(nodes.internet.status)}`">
                <span v-if="nodes.internet.status === 'pending'">...</span>
                <span v-else>{{ nodes.internet.lat.toFixed(0) }}ms</span>
            </div>
            <span class="mt-2 text-xs font-mono text-gray-400 uppercase">{{ config.internet.name }}</span>
             <span class="text-[10px] text-gray-600 font-mono">{{ config.internet.ip }}</span>
        </div>

    </div>
</template>