<script lang="ts">
  import type { KnowledgeCategory, KnowledgeItem } from '../types';
  import {
    X,
    BookOpen,
    Wrench,
    Bot,
    Boxes,
    Terminal,
    Check,
    Tag,
    Globe,
    User,
    Database,
    Image as ImageIcon
  } from '@lucide/svelte';

  export let show: boolean = false;
  export let imageSnapshot: string = '';
  export let initialTitle: string = '';
  export let initialUrl: string = '';
  export let initialAuthor: string = '';
  export let initialCategory: KnowledgeCategory = 'research';
  export let initialSummary: string = '';
  export let initialTags: string[] = [];
  export let bucketName: string = 'developer-research';

  export let onSave: (item: Partial<KnowledgeItem>) => void;
  export let onClose: () => void;

  let title = '';
  let url = '';
  let author = '';
  let category: KnowledgeCategory = 'research';
  let summary = '';
  let tags: string[] = [];
  let newTagInput = '';

  $: if (show) {
    title = initialTitle;
    url = initialUrl;
    author = initialAuthor;
    category = initialCategory;
    summary = initialSummary;
    tags = [...initialTags];
  }

  const categories: Array<{ id: KnowledgeCategory; label: string; icon: any }> = [
    { id: 'research', label: 'Research', icon: BookOpen },
    { id: 'tool', label: 'Tool', icon: Wrench },
    { id: 'agent', label: 'Agent', icon: Bot },
    { id: 'model', label: 'Model', icon: Boxes },
    { id: 'skill', label: 'Skill', icon: Terminal }
  ];

  function addTag() {
    const trimmed = newTagInput.trim().replace(/^#/, '');
    if (trimmed && !tags.includes(trimmed)) {
      tags = [...tags, trimmed];
      newTagInput = '';
    }
  }

  function removeTag(tagToRemove: string) {
    tags = tags.filter((t) => t !== tagToRemove);
  }

  function handleSave() {
    onSave({
      title: title.trim() || 'Untitled Clipping',
      url,
      author: author.trim() || 'unknown',
      category,
      summary: summary.trim(),
      tags,
      imageSnapshot,
      notes: ''
    });
  }
</script>

{#if show}
  <div
    class="modal-backdrop"
    on:click={(e) => e.target === e.currentTarget && onClose()}
    on:keydown={(e) => e.key === 'Escape' && onClose()}
    role="presentation"
  >
    <div class="modal-content" role="dialog" aria-modal="true" tabindex="-1">
      <!-- Header -->
      <div class="modal-header">
        <div class="title-group">
          <div class="icon-wrap">
            <Database size={15} strokeWidth={2} />
          </div>
          <div class="header-text">
            <h2>Clip to Knowledge Bucket</h2>
            <span class="subtext">Target: {bucketName}</span>
          </div>
        </div>
        <button class="close-btn" on:click={onClose} aria-label="Close modal">
          <X size={15} strokeWidth={2} />
        </button>
      </div>

      <!-- Body -->
      <div class="modal-body">
        <!-- Snapshot Image Preview -->
        {#if imageSnapshot}
          <div class="snapshot-preview">
            <img src={imageSnapshot} alt="Page snapshot" />
            <div class="preview-badge">
              <ImageIcon size={11} strokeWidth={2} />
              <span>Snapshot Clip</span>
            </div>
          </div>
        {/if}

        <!-- Category Selector -->
        <div class="field-group">
          <span class="field-label">Knowledge Category</span>
          <div class="category-pills">
            {#each categories as cat}
              {@const Icon = cat.icon}
              <button
                type="button"
                class="category-pill {category === cat.id ? 'selected' : ''}"
                on:click={() => (category = cat.id)}
              >
                <svelte:component this={Icon} size={13} strokeWidth={1.8} />
                <span>{cat.label}</span>
              </button>
            {/each}
          </div>
        </div>

        <!-- Title -->
        <div class="field-group">
          <label class="field-label" for="clip-title">Title</label>
          <input
            id="clip-title"
            type="text"
            class="input-text"
            bind:value={title}
            placeholder="Paper / Tool / Model name"
          />
        </div>

        <!-- URL & Author (2-column) -->
        <div class="field-row">
          <div class="field-group flex-1">
            <label class="field-label" for="clip-url">
              <Globe size={11} strokeWidth={1.8} /> Source URL
            </label>
            <input
              id="clip-url"
              type="text"
              class="input-text monospace"
              bind:value={url}
              readonly
            />
          </div>

          <div class="field-group w-36">
            <label class="field-label" for="clip-author">
              <User size={11} strokeWidth={1.8} /> Author / Handle
            </label>
            <input
              id="clip-author"
              type="text"
              class="input-text"
              bind:value={author}
              placeholder="@handle or team"
            />
          </div>
        </div>

        <!-- Extracted Summary / Insights -->
        <div class="field-group">
          <label class="field-label" for="clip-summary">Extracted Summary / Key Insights</label>
          <textarea
            id="clip-summary"
            class="input-textarea"
            rows="3"
            bind:value={summary}
            placeholder="Key developer takeaways, model details, or architecture highlights..."
          ></textarea>
        </div>

        <!-- Tags -->
        <div class="field-group">
          <span class="field-label">
            <Tag size={11} strokeWidth={1.8} /> Tags
          </span>
          <div class="tags-container">
            {#each tags as tag}
              <span class="tag-chip">
                #{tag}
                <button type="button" class="tag-remove" on:click={() => removeTag(tag)}>
                  <X size={10} strokeWidth={2} />
                </button>
              </span>
            {/each}
            <input
              type="text"
              class="tag-input"
              placeholder="Add tag..."
              bind:value={newTagInput}
              on:keydown={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault();
                  addTag();
                }
              }}
            />
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="modal-footer">
        <button type="button" class="btn-secondary" on:click={onClose}>
          Cancel
        </button>
        <button type="button" class="btn-primary" on:click={handleSave}>
          <Check size={14} strokeWidth={2.2} />
          <span>Save to Memron Bucket</span>
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 999;
    animation: fadeIn 0.15s ease-out;
  }

  .modal-content {
    background: #121215;
    border: 1px solid #27272a;
    border-radius: 0.75rem;
    width: 92%;
    max-width: 560px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 20px 40px -10px rgba(0, 0, 0, 0.6);
    overflow: hidden;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.85rem 1.1rem;
    border-bottom: 1px solid #27272a;
    background: #18181b;
  }

  .title-group {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }

  .icon-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 0.375rem;
    background: #27272a;
    color: #fafafa;
  }

  .header-text h2 {
    margin: 0;
    font-size: 0.875rem;
    font-weight: 600;
    color: #fafafa;
    letter-spacing: -0.01em;
  }

  .header-text .subtext {
    font-size: 0.6875rem;
    color: #a1a1aa;
    font-family: monospace;
  }

  .close-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
    border-radius: 0.375rem;
    border: none;
    background: none;
    color: #71717a;
    cursor: pointer;
  }

  .close-btn:hover {
    background: #27272a;
    color: #fafafa;
  }

  .modal-body {
    padding: 1.1rem;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
  }

  .snapshot-preview {
    position: relative;
    border-radius: 0.5rem;
    overflow: hidden;
    border: 1px solid #27272a;
    max-height: 140px;
    background: #09090b;
  }

  .snapshot-preview img {
    width: 100%;
    height: 140px;
    object-fit: cover;
    object-position: top;
    display: block;
  }

  .preview-badge {
    position: absolute;
    bottom: 6px;
    left: 8px;
    display: flex;
    align-items: center;
    gap: 0.35rem;
    background: rgba(9, 9, 11, 0.85);
    border: 1px solid #27272a;
    border-radius: 0.25rem;
    padding: 0.2rem 0.45rem;
    font-size: 0.625rem;
    font-weight: 500;
    color: #d4d4d8;
  }

  .field-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .field-row {
    display: flex;
    gap: 0.6rem;
  }

  .flex-1 {
    flex: 1;
  }

  .w-36 {
    width: 150px;
  }

  .field-label {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #a1a1aa;
  }

  .category-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
  }

  .category-pill {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.35rem 0.65rem;
    border-radius: 0.375rem;
    border: 1px solid #27272a;
    background: #18181b;
    color: #a1a1aa;
    font-size: 0.75rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.12s ease;
  }

  .category-pill:hover {
    border-color: #3f3f46;
    color: #fafafa;
  }

  .category-pill.selected {
    background: #27272a;
    border-color: #52525b;
    color: #fafafa;
  }

  .input-text,
  .input-textarea {
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    padding: 0.45rem 0.65rem;
    color: #fafafa;
    font-size: 0.8125rem;
    font-family: inherit;
    outline: none;
    transition: border-color 0.12s ease;
  }

  .input-text:focus,
  .input-textarea:focus {
    border-color: #52525b;
  }

  .monospace {
    font-family: monospace;
    font-size: 0.75rem;
    color: #a1a1aa;
  }

  .input-textarea {
    resize: vertical;
    min-height: 60px;
    line-height: 1.4;
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.35rem;
    background: #18181b;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    padding: 0.35rem 0.5rem;
    min-height: 34px;
  }

  .tag-chip {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    background: #27272a;
    color: #fafafa;
    border-radius: 0.25rem;
    padding: 0.15rem 0.45rem;
    font-size: 0.6875rem;
    font-family: monospace;
  }

  .tag-remove {
    display: flex;
    align-items: center;
    border: none;
    background: none;
    color: #a1a1aa;
    cursor: pointer;
    padding: 0;
  }

  .tag-remove:hover {
    color: #fafafa;
  }

  .tag-input {
    border: none;
    background: none;
    color: #fafafa;
    font-size: 0.75rem;
    outline: none;
    flex: 1;
    min-width: 80px;
  }

  .modal-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 0.5rem;
    padding: 0.8rem 1.1rem;
    border-top: 1px solid #27272a;
    background: #18181b;
  }

  .btn-secondary {
    padding: 0.45rem 0.85rem;
    border-radius: 0.375rem;
    border: 1px solid #27272a;
    background: none;
    color: #a1a1aa;
    font-size: 0.75rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.12s ease;
  }

  .btn-secondary:hover {
    background: #27272a;
    color: #fafafa;
  }

  .btn-primary {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.45rem 0.95rem;
    border-radius: 0.375rem;
    border: none;
    background: #fafafa;
    color: #09090b;
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.12s ease;
  }

  .btn-primary:hover {
    background: #e4e4e7;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @media (max-width: 520px) {
    .modal-content {
      width: 95%;
      max-height: 94vh;
    }

    .modal-body {
      padding: 0.85rem;
      gap: 0.75rem;
    }

    .field-row {
      flex-direction: column;
      gap: 0.75rem;
    }

    .w-36 {
      width: 100%;
    }

    .modal-footer {
      padding: 0.65rem 0.85rem;
    }
  }
</style>
