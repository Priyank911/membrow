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
  export let processing: boolean = false;
  export let extracted: boolean = false;

  export let onSave: (item: Partial<KnowledgeItem>) => void;
  export let onExtract: () => void;
  export let onClose: () => void;

  let title = '';
  let url = '';
  let author = '';
  let category: KnowledgeCategory = 'research';
  let summary = '';
  let tags: string[] = [];
  let newTagInput = '';
  let zoom = 1;

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

  function adjustZoom(amount: number) {
    zoom = Math.min(2, Math.max(0.5, Number((zoom + amount).toFixed(2))));
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
        <button class="close-btn" on:click={onClose} aria-label="Close modal" disabled={processing}>
          <X size={15} strokeWidth={2} />
        </button>
      </div>

      <!-- Snapshot preview -->
      <div class="capture-pane">
          {#if imageSnapshot}
            <button
              type="button"
              class="capture-stage"
              on:click={onExtract}
              disabled={processing || extracted}
              aria-label={extracted ? 'Snapshot extracted' : 'Extract knowledge from snapshot'}
            >
              <img
                src={imageSnapshot}
                alt="Page snapshot"
                style={`transform: scale(${zoom});`}
              />
            </button>
            <div class="capture-toolbar">
              <span class="capture-label"><ImageIcon size={12} /> CAPTURE VIEW</span>
              <div class="zoom-controls">
                <button type="button" on:click={() => adjustZoom(-0.1)} aria-label="Zoom out">−</button>
                <span>{Math.round(zoom * 100)}%</span>
                <button type="button" on:click={() => adjustZoom(0.1)} aria-label="Zoom in">+</button>
              </div>
            </div>
          {/if}
      </div>

      <!-- Metadata form -->
      <div class="modal-body">
        <div class="form-heading">
          <span class="eyebrow">KNOWLEDGE ENTRY</span>
          <span class="form-hint">Review the extracted signal before saving.</span>
        </div>
        {#if processing}
          <div class="processing-banner" role="status">
            <span class="processing-spinner"></span>
            <div>
              <strong>Processing snapshot...</strong>
              <span>Extracting knowledge and confirming remote Memron storage.</span>
            </div>
          </div>
        {/if}

        <div class="field-group">
          <span class="field-label">Knowledge Category</span>
          <div class="category-pills">
            {#each categories as cat}
              {@const Icon = cat.icon}
              <button
                type="button"
                class="category-pill {category === cat.id ? 'selected' : ''}"
                on:click={() => (category = cat.id)}
                disabled={processing || !extracted}
              >
                <svelte:component this={Icon} size={13} strokeWidth={1.8} />
                <span>{cat.label}</span>
              </button>
            {/each}
          </div>
        </div>

        <div class="field-group">
          <label class="field-label" for="clip-title">Title</label>
          <input
            id="clip-title"
            type="text"
            class="input-text"
            bind:value={title}
            placeholder="Paper / Tool / Model name"
            disabled={processing || !extracted}
          />
        </div>

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
              disabled={processing || !extracted}
            />
          </div>
        </div>

        <div class="field-group">
          <label class="field-label" for="clip-summary">Extracted Summary / Key Insights</label>
          <textarea
            id="clip-summary"
            class="input-textarea"
            rows="3"
            bind:value={summary}
            placeholder="Key developer takeaways, model details, or architecture highlights..."
            disabled={processing || !extracted}
          ></textarea>
        </div>

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
              disabled={processing || !extracted}
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
        <button type="button" class="btn-secondary" on:click={onClose} disabled={processing}>
          Cancel
        </button>
        <button type="button" class="btn-primary" on:click={handleSave} disabled={processing || !extracted}>
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
    background: rgba(0, 0, 0, 0.86);
    backdrop-filter: blur(10px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 999;
    animation: fadeIn 0.15s ease-out;
  }

  .modal-content {
    background: #0b0b0c;
    border: 1px solid #29292d;
    border-radius: 12px;
    width: min(1180px, calc(100vw - 48px));
    height: min(650px, calc(100vh - 40px));
    display: grid;
    grid-template-columns: minmax(0, 1.05fr) minmax(380px, 0.95fr);
    grid-template-rows: auto minmax(0, 1fr) auto;
    box-shadow: 0 30px 100px rgba(0, 0, 0, 0.75);
    overflow: hidden;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    grid-column: 1 / -1;
    padding: 0.85rem 1.1rem;
    border-bottom: 1px solid #242424;
    background: #101012;
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
    width: 32px;
    height: 32px;
    border-radius: 0.55rem;
    background: linear-gradient(145deg, #25252b, #17171a);
    border: 1px solid #3a3a42;
    color: #f4f4f5;
  }

  .header-text h2 {
    margin: 0;
    font-size: 0.95rem;
    font-weight: 700;
    color: #fafafa;
    letter-spacing: -0.01em;
  }

  .header-text .subtext {
    font-size: 0.6875rem;
    color: #71717a;
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
    grid-column: 2;
    grid-row: 2;
    padding: 1.1rem 1.2rem;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    min-width: 0;
    border-left: 1px solid #2b2b30;
  }

  .form-heading {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.75rem;
    padding-bottom: 0.65rem;
    border-bottom: 1px solid #202024;
  }

  .eyebrow {
    color: #d4d4d8;
    font-size: 0.62rem;
    font-weight: 700;
    letter-spacing: 0.1em;
  }

  .form-hint {
    color: #52525b;
    font-size: 0.65rem;
  }

  .processing-banner {
    display: flex;
    align-items: center;
    gap: 0.7rem;
    padding: 0.7rem 0.8rem;
    border: 1px solid #3f3f46;
    border-radius: 0.5rem;
    background: #18181b;
    color: #e4e4e7;
    position: fixed;
    right: 1.25rem;
    bottom: 1.25rem;
    z-index: 4;
    min-width: 290px;
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.65);
  }

  .processing-banner strong,
  .processing-banner span {
    display: block;
  }

  .processing-banner strong {
    font-size: 0.75rem;
  }

  .processing-banner div span {
    margin-top: 0.15rem;
    color: #a1a1aa;
    font-size: 0.6875rem;
  }

  .processing-spinner {
    width: 1rem;
    height: 1rem;
    border: 2px solid #52525b;
    border-top-color: #fafafa;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .capture-pane {
    grid-column: 1;
    grid-row: 2 / 4;
    min-width: 0;
    min-height: 0;
    display: flex;
    flex-direction: column;
    padding: 1.1rem;
    background: #070708;
    min-width: 0;
    border-right: 1px solid #3a3a42;
    box-shadow: 1px 0 0 #101012;
  }

  .capture-stage {
    appearance: none;
    width: 100%;
    border: 1px solid #1e1e1e;
    border-radius: 10px;
    color: inherit;
    padding: 0;
    flex: 1;
    min-height: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: auto;
    background: #0d0d0f;
    box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.02);
    cursor: pointer;
  }

  .capture-stage:hover:not(:disabled) {
    border-color: #444;
  }

  .capture-stage:disabled {
    cursor: default;
  }

  .capture-stage img {
    width: auto;
    height: auto;
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    transform-origin: center center;
    transition: transform 0.16s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .capture-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.65rem 0.1rem 0;
    color: #777;
    font-size: 0.62rem;
    letter-spacing: 0.12em;
  }

  .capture-label {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .zoom-controls {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #bbb;
    letter-spacing: 0;
  }

  .zoom-controls button {
    width: 24px;
    height: 24px;
    border: 1px solid #2b2b2b;
    border-radius: 6px;
    background: #111;
    color: #eee;
    cursor: pointer;
    font-size: 1rem;
    line-height: 1;
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
    grid-column: 2;
    grid-row: 3;
    padding: 0.7rem 1.2rem 0.85rem;
    background: #0b0b0c;
    border-left: 1px solid #2b2b30;
    border-top: 1px solid #202024;
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
      height: 94vh;
      grid-template-columns: 1fr;
      grid-template-rows: min-content minmax(210px, 0.8fr) minmax(0, 1.4fr) auto;
    }

    .capture-pane {
      grid-column: 1;
      grid-row: 2;
      padding: 0.85rem;
    }

    .modal-body {
      grid-column: 1;
      grid-row: 3;
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
      grid-column: 1;
      grid-row: 4;
      padding: 0.65rem 0.85rem;
      border-left: 0;
    }
  }
</style>
