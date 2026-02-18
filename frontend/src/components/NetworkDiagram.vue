<script setup lang="ts">
import { ref, onMounted } from 'vue';

// Estado dos nós
const nodes = ref({
    local: { name: "You (Local)", status: "ok", lat: 0 },
    gateway: { name: "Gateway", status: "pending", lat: 0, ip: "192.168.1.1" },
    internet: { name: "Internet (8.8.8.8)", status: "pending", lat: 0, ip: "8.8.8.8" }
});

const updateNode = (id: 'gateway' | 'internet', latency: number, loss: boolean) => {
    const node = nodes.value[id];
    node.lat = latency / 1000; // Converte micro -> mili
    
    if (loss) {
        node.status = "error";
    } else if (node.lat > 100) {
        node.status = "warn";
    } else {
        node.status = "ok";
    }
};

onMounted(() => {
    // Escuta eventos do backend
    // @ts-ignore
    window.runtime.EventsOn("ping:data", (data: any) => {
        // Identifica qual nó atualizar baseado no IP ou ID
        if (data.hostId === "gateway") {
            updateNode("gateway", data.latency, data.loss);
        } else if (data.hostId === "google") { // ID definido no main.go
            updateNode("internet", data.latency, data.loss);
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
            <span class="mt-2 text-xs font-mono text-gray-400">LOCALHOST</span>
        </div>

        <div class="flex-1 h-1 bg-gray-800 mx-2 relative overflow-hidden">
            <div class="absolute inset-0 bg-cyan-500/20 animate-pulse"></div>
            <div class="absolute -top-6 left-1/2 -translate-x-1/2 text-xs font-mono text-cyan-400">
                LAN
            </div>
        </div>

        <div class="flex flex-col items-center z-10 transition-all duration-300">
            <div :class="`w-16 h-16 rounded-full flex items-center justify-center border-2 border-transparent text-black font-bold transition-all duration-300 ${getStatusColor(nodes.gateway.status)}`">
                <span v-if="nodes.gateway.status === 'pending'">...</span>
                <span v-else>{{ nodes.gateway.lat.toFixed(1) }}ms</span>
            </div>
            <span class="mt-2 text-xs font-mono text-gray-400">GATEWAY</span>
        </div>

        <div class="flex-1 h-1 bg-gray-800 mx-2 relative">
             <div :class="`w-full h-full transition-colors duration-500 ${nodes.gateway.lat > 50 ? 'bg-yellow-500/50' : 'bg-cyan-500/20'}`"></div>
        </div>

        <div class="flex flex-col items-center z-10">
            <div :class="`w-16 h-16 rounded-full flex items-center justify-center border-2 border-transparent text-black font-bold transition-all duration-300 ${getStatusColor(nodes.internet.status)}`">
                <span v-if="nodes.internet.status === 'pending'">...</span>
                <span v-else>{{ nodes.internet.lat.toFixed(1) }}ms</span>
            </div>
            <span class="mt-2 text-xs font-mono text-gray-400">INTERNET</span>
        </div>

    </div>
</template>