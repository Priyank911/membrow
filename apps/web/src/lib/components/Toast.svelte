<script lang="ts">
  import type { ToastMessage } from '../types';
  import { Check, AlertCircle, Info, X } from '@lucide/svelte';

  export let toasts: ToastMessage[] = [];

  function remove(id: string) {
    toasts = toasts.filter((t) => t.id !== id);
  }
</script>

{#if toasts.length > 0}
  <div class="toast-container" role="region" aria-label="Notifications">
    {#each toasts as toast (toast.id)}
      <div class="toast {toast.type || 'info'}">
        <span class="icon">
          {#if toast.type === 'success'}
            <Check size={14} strokeWidth={2} />
          {:else if toast.type === 'error'}
            <AlertCircle size={14} strokeWidth={2} />
          {:else}
            <Info size={14} strokeWidth={2} />
          {/if}
        </span>
        <span class="text">{toast.text}</span>
        <button class="close-btn" on:click={() => remove(toast.id)} aria-label="Dismiss">
          <X size={12} />
        </button>
      </div>
    {/each}
  </div>
{/if}

<style>
  .toast-container {
    position: fixed;
    bottom: 1.5rem;
    right: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    z-index: 1000;
    pointer-events: none;
  }

  .toast {
    pointer-events: auto;
    display: flex;
    align-items: center;
    gap: 0.6rem;
    background: #18181b;
    border: 1px solid #27272a;
    color: #fafafa;
    padding: 0.5rem 0.8rem;
    border-radius: 0.5rem;
    font-size: 0.8125rem;
    font-weight: 500;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5);
    animation: slideUp 0.18s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .toast.success {
    border-color: #22c55e33;
    color: #f0fdf4;
  }

  .toast.success .icon {
    color: #22c55e;
  }

  .toast.error {
    border-color: #ef444433;
    color: #fef2f2;
  }

  .toast.error .icon {
    color: #ef4444;
  }

  .toast.info .icon {
    color: #a1a1aa;
  }

  .icon {
    display: flex;
    align-items: center;
  }

  .text {
    letter-spacing: -0.01em;
  }

  .close-btn {
    background: none;
    border: none;
    color: #71717a;
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 0.1rem;
    margin-left: 0.25rem;
  }

  .close-btn:hover {
    color: #fafafa;
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
