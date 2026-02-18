<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { UpdateConfig, GetConfig } from '../../wailsjs/go/main/App'; // Ajustado para os métodos do ConfigManager

const retentionDays = ref(7);
const isSaving = ref(false);
const showSuccess = ref(false);

onMounted(async () => {
    try {
        // Busca a configuração atual do settings.json
        const config = await GetConfig();
        if (config && config.retention_days) {
            retentionDays.value = config.retention_days;
        }
    } catch (err) {
        console.error("Erro ao carregar configurações:", err);
    }
});

const handleSave = async () => {
    isSaving.value = true;
    try {
        // 1. Busca a instância da classe AppConfig vinda do backend
        const currentConfig = await GetConfig();
        
        // 2. Modifica a propriedade diretamente na instância
        // Isso preserva os métodos internos (como convertValues) que o Wails injeta
        currentConfig.retention_days = Number(retentionDays.value);

        // 3. Envia a instância modificada de volta
        await UpdateConfig(currentConfig);
        
        showSuccess.value = true;
        setTimeout(() => showSuccess.value = false, 3000);
    } catch (err) {
        alert("Erro ao salvar no settings.json: " + err);
    } finally {
        isSaving.value = false;
    }
};
</script>

<template>
    <div class="max-w-2xl mx-auto p-6 bg-gray-900/40 border border-gray-800 rounded-xl">
        <h2 class="text-green-400 font-mono mb-6 uppercase text-sm tracking-widest">// APP SETTINGS (FILE PERSISTENCE)</h2>
        
        <div class="space-y-6">
            <div class="p-6 bg-black/40 border border-gray-800 rounded-lg group hover:border-green-500/30 transition-all">
                <label class="block text-[10px] font-mono text-gray-500 uppercase mb-4 tracking-widest">
                    Retention Policy (Days)
                </label>
                
                <div class="flex items-center gap-6">
                    <input type="number" v-model="retentionDays" min="1" max="365"
                        class="bg-black border border-gray-700 p-3 rounded text-green-500 font-mono w-24 outline-none focus:border-green-500 text-lg" />
                    
                    <div class="text-xs font-mono leading-relaxed">
                        <p class="text-gray-300">Auto-cleanup is <span class="text-green-500">ENABLED</span></p>
                        <p class="text-gray-600">Os logs serão removidos do arquivo e banco após este período.</p>
                    </div>
                </div>
            </div>

            <button @click="handleSave" :disabled="isSaving"
                class="w-full bg-green-600 hover:bg-green-500 disabled:bg-gray-800 text-black font-bold py-4 rounded-lg transition-all font-mono text-xs uppercase tracking-widest">
                {{ isSaving ? 'Committing to JSON...' : 'Save Settings' }}
            </button>

            <p v-if="showSuccess" class="text-center text-green-500 font-mono text-[10px] animate-pulse">
                ✓ SETTINGS.JSON UPDATED SUCCESSFULLY
            </p>
        </div>
    </div>
</template>