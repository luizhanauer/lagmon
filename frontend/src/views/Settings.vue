<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { SaveSetting, GetSetting } from '../../wailsjs/go/main/App';

const retentionDays = ref(7);
const isSaving = ref(false);
const showSuccess = ref(false);

onMounted(async () => {
    const saved = await GetSetting('retention_days');
    if (saved) retentionDays.value = parseInt(saved);
});

const handleSave = async () => {
    isSaving.value = true;
    try {
        await SaveSetting('retention_days', retentionDays.value.toString());
        showSuccess.value = true;
        setTimeout(() => showSuccess.value = false, 3000);
    } catch (err) {
        alert("Erro ao salvar no banco: " + err);
    } finally {
        isSaving.value = false;
    }
};
</script>

<template>
    <div class="max-w-2xl mx-auto p-6 bg-gray-900/40 border border-gray-800 rounded-xl">
        <h2 class="text-green-400 font-mono mb-6 uppercase text-sm tracking-widest">// DATABASE SETTINGS</h2>
        
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
                        <p class="text-gray-600">Registros mais antigos que {{ retentionDays }} dias serão expurgados do SQLite.</p>
                    </div>
                </div>
            </div>

            <button @click="handleSave" :disabled="isSaving"
                class="w-full bg-green-600 hover:bg-green-500 disabled:bg-gray-800 text-black font-bold py-4 rounded-lg transition-all font-mono text-xs uppercase tracking-widest">
                {{ isSaving ? 'Committing to Disk...' : (showSuccess ? '✓ Settings Applied' : 'Save Configuration') }}
            </button>
        </div>
    </div>
</template>