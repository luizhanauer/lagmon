<script setup lang="ts">
import { onMounted, ref } from 'vue';
import uPlot from 'uplot';
import 'uplot/dist/uPlot.min.css';
// import { watch } from 'vue';

const props = defineProps<{
    targetId: string;
    targetName: string;
    color: string;
}>();

// watch(() => props.color, (newColor) => {
//     if (uplotInst) {
//         // Atualiza a cor da série de latência (índice 1) sem destruir o gráfico
//         uplotInst.series[1].stroke = newColor;
//         uplotInst.series[1].fill = newColor + "15";
//         uplotInst.redraw();
//     }
// });

const chartRef = ref<HTMLElement>();
let uplotInst: uPlot | null = null;
let data: any[] = [[], [], []];

// Função auxiliar para formatar Data
const fmtDate = (u: uPlot, val: number) => {
    if (!val) return "-";
    const d = new Date(val * 1000);
    const pad = (n: number) => n.toString().padStart(2, '0');
    return `${pad(d.getDate())}/${pad(d.getMonth() + 1)}/${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
};

const initChart = () => {
    if (!chartRef.value) return;

    const opts: uPlot.Options = {
        width: chartRef.value.clientWidth,
        height: 220,
        class: "bg-transparent",
        
        cursor: {
            show: true,
            drag: { x: false, y: false },
            points: { show: false }
        },
        legend: {
            show: true,
        },
        
        scales: {
            x: { time: true } // Garante que o uPlot trate como tempo
        },

        series: [
            { 
                label: "Data",
                // AQUI ESTÁ A MÁGICA: Formatador customizado para a Legenda
                value: fmtDate 
            }, 
            {
                label: "Latency",
                stroke: props.color,
                width: 2,
                fill: props.color + "15",
                points: { show: false },
                spanGaps: false,
                value: (u, v) => v == null ? "TIMEOUT" : v.toFixed(1) + "ms"            },
            {
                label: "Jitter",
                stroke: "#ffffff",
                width: 1,
                dash: [3, 3],
                points: { show: false },
                value: (u, v) => v == null ? "-" : v.toFixed(1) + "ms"
            }
        ],
        axes: [
            { 
                show: true, // Mostra o eixo X
                stroke: "#555",
                font: "10px monospace",
                // Formata os "ticks" do eixo X (embaixo do gráfico)
                values: (u, vals) => vals.map(v => {
                    const d = new Date(v * 1000);
                    const pad = (n: number) => n.toString().padStart(2, '0');
                    return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
                })
            }, 
            { 
                stroke: "#555", 
                gap: 5, 
                size: 35,
                grid: { stroke: "#333", width: 1 },
                values: (u, v) => v.map(i => i.toFixed(0)) 
            }
        ]
    };

    uplotInst = new uPlot(opts, data as any, chartRef.value);
};

const handleEvent = (payload: any) => {
    if (payload.hostId !== props.targetId) return;
    if (!uplotInst) return;

    const now = Math.floor(Date.now() / 1000);
    
    // Tratamento robusto para Packet Loss
    const isLoss = payload.loss === true;
    const lat = isLoss ? null : payload.latency / 1000;
    const jit = isLoss ? null : payload.jitter / 1000;

    data[0].push(now);
    data[1].push(lat);
    data[2].push(jit);

    if (data[0].length > 60) {
        data.forEach(ch => ch.shift());
    }
    uplotInst.setData(data as any);
};

onMounted(() => {
    setTimeout(() => {
        initChart();
        // @ts-ignore
        window.runtime.EventsOn("ping:data", handleEvent);
    }, 100);

    new ResizeObserver(() => {
        if(uplotInst && chartRef.value) {
            uplotInst.setSize({ width: chartRef.value.clientWidth, height: 220 });
        }
    }).observe(chartRef.value!);
});
</script>

<template>
    <div class="w-full h-full relative group">
        <component is="style">
            .u-legend { 
                font-family: 'Fira Code', monospace; 
                font-size: 11px; 
                color: #ccc; 
                background: rgba(0,0,0,0.85); 
                backdrop-filter: blur(4px);
                position: absolute; 
                top: 0; 
                right: 0; 
                padding: 6px 10px; 
                border-bottom-left-radius: 8px; 
                border-left: 1px solid #333;
                border-bottom: 1px solid #333;
                z-index: 10; 
                display: none; /* Esconde por padrão, o uPlot mostra no hover */
            }
            /* Mostra a legenda apenas quando o gráfico tem hover (classe interna do uPlot) */
            .uplot:hover .u-legend { display: block; }

            .u-legend tr { display: block; margin-bottom: 2px; }
            .u-series th { vertical-align: top; text-align: left; padding-right: 8px; opacity: 0.7; }
            .u-series td { font-weight: bold; text-align: right; color: #fff; }
        </component>

        <div ref="chartRef" class="w-full h-full"></div>
        
        <div class="absolute top-2 left-2 text-[10px] font-bold uppercase tracking-widest px-2 py-1 rounded bg-black/60 backdrop-blur-md border border-gray-700/50 pointer-events-none" :style="{ color: props.color }">
            {{ props.targetName }}
        </div>
    </div>
</template>