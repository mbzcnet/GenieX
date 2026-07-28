package com.geniex.sdk.bean

/**
 * Model configuration corresponding to the native `geniex_ModelConfig` struct.
 */
data class ModelConfig(
    /**
     * Text context size for llama_cpp (default 8192). Use 0 for the native
     * plugin default. For qairt the JNI forces 0 — context length is fixed by
     * the AI Hub bundle.
     */
    var nCtx: Int = 8192,

    /** Number of threads used for text generation */
    var nThreads: Int = 8,

    /** Number of threads used for batch processing */
    var nThreadsBatch: Int = 8,

    /** Maximum logical batch size submitted to llama_decode */
    var nBatch: Int = 2048,

    /** Maximum physical batch size supported by backend */
    var nUBatch: Int = 512,

    /** Maximum number of distinct sequences (states) */
    var nSeqMax: Int = 1,

    /**
     * Number of layers to offload to GPU / NPU; -1 (the default) means all
     * layers, which llama.cpp interprets natively. The JNI layer forces 0
     * when [InputPluginBase.compute_unit] is [ComputeUnitValue.CPU] (and for
     * qairt, which ignores it). For [ComputeUnitValue.GPU] / [ComputeUnitValue.NPU]
     * / [ComputeUnitValue.HYBRID] the caller's value passes through.
     */
    var nGpuLayers: Int = -1,

    /**
     * llama_cpp-only KV cache element type for K ("f16", "q8_0", "q4_0", …).
     * Empty = auto (q8_0 when [nCtx] >= 8192). Ignored by qairt.
     */
    var cacheTypeK: String = "",

    /**
     * llama_cpp-only KV cache element type for V. Empty = auto. Ignored by qairt.
     */
    var cacheTypeV: String = "",

    /** Path to the chat template file (optional) */
    val chat_template_path: String = "",

    /** Content of the chat template file (optional) */
    val chat_template_content: String = "",

    /** Maximum number of tokens to generate */
    val max_tokens: Int = 2048,

    /** Enable "thinking" mode for more detailed reasoning */
    val enable_thinking: Boolean = false,

    val verbose: Boolean = false,
)
