<script setup lang="ts">
import { computed } from 'vue';
import { store } from '../store';

// Filtra apenas os alvos que o usuário selecionou para o diagrama
const diagramTargets = computed(() => store.targets.filter(t => t.showInDiagram));

/**
 * Define o estilo visual baseado no status do ping.
 * Mantém uniformidade e destaca anomalias com efeitos neon.
 */
const getStatusClasses = (stats: any, active: boolean) => {
    if (!active) return 'border-gray-800 bg-gray-900/20 text-gray-600 opacity-50';
    
    if (stats.loss) {
        return 'border-red-500 bg-red-500/10 text-red-500 shadow-[0_0_15px_rgba(239,68,68,0.2)] animate-pulse';
    }
    
    if (stats.latency > 150) {
        return 'border-yellow-500 bg-yellow-500/10 text-yellow-500 shadow-[0_0_15px_rgba(250,204,21,0.2)]';
    }
    
    return 'border-cyan-500 bg-cyan-500/5 text-cyan-400 shadow-[0_0_10px_rgba(6,182,212,0.1)]';
};
</script>

<template>
    <div class="w-full space-y-4 py-2">
        <div class="flex items-center gap-2 px-1">
            <div class="w-1 h-3 bg-green-500"></div>
            <h2 class="text-[10px] font-mono text-gray-500 uppercase tracking-widest">
                Monitoring Nodes
            </h2>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
            <div 
                v-for="target in diagramTargets" 
                :key="target.id" 
                :class="[
                    'relative flex flex-col items-center justify-center p-4 rounded-xl border-2 transition-all duration-500 aspect-square min-w-[120px]',
                    getStatusClasses(target.stats, target.active)
                ]"
            >
                <span v-if="target.isGateway" 
                    class="absolute top-2 right-2 text-[8px] font-bold px-1 border border-current rounded opacity-60">
                    GW
                </span>

                <div class="text-xl mb-2 opacity-80">
                    {{ target.isGateway ? '📡' : '🌐' }}
                </div>
                
                <div class="text-center overflow-hidden w-full">
                    <div class="font-bold text-[10px] uppercase tracking-tighter truncate">
                        {{ target.name }}
                    </div>
                    <div class="font-mono text-[9px] opacity-40 truncate">
                        {{ target.ip }}
                    </div>
                </div>

                <div class="mt-2 font-mono text-lg font-black tracking-tight leading-none">
                    {{ target.stats.loss ? 'LOSS' : `${target.stats.latency.toFixed(0)}ms` }}
                </div>
            </div>

            <div v-if="diagramTargets.length === 0" 
                class="col-span-full py-10 border-2 border-dashed border-gray-800 rounded-xl flex flex-col items-center justify-center bg-gray-900/10">
                <p class="text-gray-600 font-mono text-[10px] uppercase tracking-[0.2em]">
                    No paths active in diagram mode
                </p>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Efeito de scanline sutil para cards em estado de erro */
.animate-pulse {
    background: linear-gradient(
        0deg, 
        rgba(239, 68, 68, 0.05) 0%, 
        rgba(239, 68, 68, 0.15) 50%, 
        rgba(239, 68, 68, 0.05) 100%
    );
    background-size: 100% 4px;
}
</style>