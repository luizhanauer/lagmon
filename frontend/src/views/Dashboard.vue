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
            <NetworkDiagram />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-6 pb-10">
            
            <div v-for="target in visibleTargets" :key="target.id" 
                 :class="[
                    'bg-gray-900/40 border rounded-xl overflow-hidden shadow-lg transition-all relative flex flex-col h-[260px]',
                    target.stats.loss ? 'border-red-900/50 shadow-red-900/10' : 'border-gray-800 hover:border-gray-700'
                 ]">
                
                <div v-if="!target.active" class="absolute inset-0 z-20 bg-black/70 flex items-center justify-center backdrop-blur-sm">
                    <span class="text-gray-500 font-mono text-sm tracking-widest">[ MONITORING PAUSED ]</span>
                </div>

                <div v-if="target.stats.loss && target.active" class="absolute top-2 right-2 z-10">
                    <span class="flex h-2 w-2">
                        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                        <span class="relative inline-flex rounded-full h-2 w-2 bg-red-500"></span>
                    </span>
                </div>

                <div class="flex-1 w-full relative">
                    <PingChart 
                        :targetId="target.id" 
                        :targetName="target.name" 
                        :color="target.stats.loss ? '#ef4444' : target.color" 
                    />
                </div>

                <div :class="[
                    'h-10 px-4 flex justify-between items-center text-xs font-mono border-t shrink-0 transition-colors',
                    target.stats.loss ? 'bg-red-500/10 border-red-900/30' : 'bg-black/40 border-gray-800'
                ]">
                    <span class="text-gray-500">{{ target.ip }}</span>
                    
                    <div class="flex gap-4">
                        <div class="flex items-center gap-1">
                            <span class="text-[10px] text-gray-600 uppercase">Lat:</span>
                            <b :class="target.stats.loss ? 'text-red-500 animate-pulse' : 'text-white text-sm'">
                                {{ target.stats.loss ? 'TIMEOUT' : `${target.stats.latency.toFixed(0)}ms` }}
                            </b>
                        </div>
                        <div class="flex items-center gap-1">
                            <span class="text-[10px] text-gray-600 uppercase">Jit:</span>
                            <b class="text-gray-400 text-sm">{{ target.stats.jitter.toFixed(0) }}ms</b>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="visibleTargets.length === 0" class="col-span-full py-12 text-center border-2 border-dashed border-gray-800 rounded-xl bg-gray-900/20">
                <p class="text-gray-600 font-mono text-sm">NO GRAPHS SELECTED IN TARGETS TAB</p>
            </div>

        </div>
    </div>
</template>