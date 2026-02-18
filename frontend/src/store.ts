import { reactive } from 'vue';

export interface HostDef {
    id: string;
    ip: string;
    name: string;
    color: string;
    showGraph: boolean;
    active: boolean;
    stats: { latency: number; jitter: number; loss: boolean };
}

const NEON_PALETTE = ["#00E5FF", "#D946EF", "#FACC15", "#4ADE80", "#FF5722", "#FFFFFF"];

export const store = reactive({
    targets: [] as HostDef[],
    
    addTarget(id: string, ip: string, name: string, active: boolean = true) {
        if (this.targets.find(t => t.id === id)) return;

        const color = NEON_PALETTE[this.targets.length % NEON_PALETTE.length];
        this.targets.push({
            id, ip, name, color,
            showGraph: true, // Padrão: mostrar gráfico
            active: active,
            stats: { latency: 0, jitter: 0, loss: false }
        });
    },

    removeTarget(id: string) {
        this.targets = this.targets.filter(t => t.id !== id);
    },

    updateStats(hostId: string, latency: number, jitter: number, loss: boolean) {
        const target = this.targets.find(t => t.id === hostId);
        if (target) {
            target.stats = { latency, jitter, loss };
        }
    }
});