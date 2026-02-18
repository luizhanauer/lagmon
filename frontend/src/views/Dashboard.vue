<script setup lang="ts">
import { computed } from 'vue';
import { store } from '../store';
import PingChart from '../components/PingChart.vue';
import NetworkDiagram from '../components/NetworkDiagram.vue';

const visibleTargets = computed(() => store.targets.filter(t => t.showGraph));
</script>

<template>
    <div class="flex flex-col gap-6 w-full">
        
        <div class="w-full">
            <h2 class="text-xs font-mono text-gray-500 mb-2 uppercase tracking-wider">// Network Topology</h2>
            <NetworkDiagram />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-6 pb-10">
            
            <div v-for="target in visibleTargets" :key="target.id" 
                 class="bg-gray-900/40 border border-gray-800 rounded-xl overflow-hidden shadow-lg transition-all hover:border-gray-700 relative flex flex-col h-[260px]">
                
                <div v-if="!target.active" class="absolute inset-0 z-20 bg-black/70 flex items-center justify-center backdrop-blur-sm">
                    <span class="text-gray-500 font-mono text-sm">[ PAUSED ]</span>
                </div>

                <div class="flex-1 w-full relative">
                    <PingChart 
                        :targetId="target.id" 
                        :targetName="target.name" 
                        :color="target.color" 
                    />
                </div>

                <div class="h-10 px-4 bg-black/40 flex justify-between items-center text-xs font-mono border-t border-gray-800 shrink-0">
                    <span class="text-gray-500">{{ target.ip }}</span>
                    
                    <div class="flex gap-4">
                        <span :class="target.stats.loss ? 'text-red-500' : 'text-gray-300'">
                            LAT: <b :class="target.stats.loss ? 'text-red-500' : 'text-white text-sm'">{{ target.stats.loss ? 'ERR' : target.stats.latency.toFixed(0) }}</b>ms
                        </span>
                        <span class="text-gray-500">
                            JIT: <b class="text-gray-400 text-sm">{{ target.stats.jitter.toFixed(0) }}</b>ms
                        </span>
                    </div>
                </div>
            </div>

            <div v-if="visibleTargets.length === 0" class="col-span-full py-12 text-center border-2 border-dashed border-gray-800 rounded-xl bg-gray-900/20">
                <p class="text-gray-600 font-mono text-sm">NO GRAPHS SELECTED</p>
            </div>

        </div>
    </div>
</template>