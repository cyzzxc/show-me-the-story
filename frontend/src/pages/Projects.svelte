<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';
  import { currentProject, projects, addToast, showConfirm, taskRunning, progress, config, settings, chatSessions, currentChatSession } from '../lib/stores.js';

  let newProjectName = '';
  let creating = false;

  const phaseLabels = { outline: '大纲', writing: '写作' };

  onMount(loadProjects);

  async function loadProjects() {
    try {
      const list = await api('GET', '/api/projects');
      projects.set(Array.isArray(list) ? list : []);
    } catch (e) {
      projects.set([]);
    }
  }

  async function selectProject(name) {
    try {
      await api('POST', '/api/projects/select', { name });
      currentProject.set(name);
      // Reload all project data
      try { progress.set(await api('GET', '/api/progress')); } catch (e) {}
      try { config.set(await api('GET', '/api/config')); } catch (e) {}
      try { settings.set(await api('GET', '/api/settings')); } catch (e) {}
      try { chatSessions.set(await api('GET', '/api/chat/sessions')); } catch (e) {}
      currentChatSession.set(null);
      addToast('已切换到项目: ' + name, 'success');
    } catch (e) {
      addToast(e.message, 'error');
    }
  }

  async function createProject() {
    const name = newProjectName.trim();
    if (!name) {
      addToast('请输入项目名称', 'error');
      return;
    }
    creating = true;
    try {
      await api('POST', '/api/projects', { name });
      newProjectName = '';
      await loadProjects();
      await selectProject(name);
    } catch (e) {
      addToast(e.message, 'error');
    } finally {
      creating = false;
    }
  }

  async function deleteProject(name) {
    showConfirm(`确认删除项目「${name}」？此操作不可恢复！`, async () => {
      try {
        await api('DELETE', '/api/projects/' + encodeURIComponent(name));
        await loadProjects();
        addToast('项目已删除', 'success');
      } catch (e) {
        addToast(e.message, 'error');
      }
    });
  }

  function handleKeydown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      createProject();
    }
  }
</script>

<div class="flex items-center justify-center min-h-[60vh]">
  <div class="w-full max-w-lg space-y-6">
    <!-- Title -->
    <div class="text-center">
      <div class="text-5xl mb-4">📚</div>
      <h2 class="text-2xl font-bold mb-1">选择故事项目</h2>
      <p class="text-sm text-base-content/50">选择已有项目或创建新项目开始创作</p>
    </div>

    <!-- Create new project -->
    <div class="card bg-base-200 shadow-sm">
      <div class="card-body p-4">
        <h3 class="card-title text-sm">新建项目</h3>
        <div class="flex gap-2">
          <input
            type="text"
            class="input input-sm flex-1"
            bind:value={newProjectName}
            placeholder="输入项目名称..."
            on:keydown={handleKeydown}
            disabled={creating}
          />
          <button
            class="btn btn-primary btn-sm"
            on:click={createProject}
            disabled={creating || !newProjectName.trim()}
          >
            {#if creating}
              <span class="loading loading-spinner loading-xs"></span>
            {:else}
              创建
            {/if}
          </button>
        </div>
      </div>
    </div>

    <!-- Project list -->
    <div class="card bg-base-200 shadow-sm">
      <div class="card-body p-4">
        <h3 class="card-title text-sm">已有项目 <span class="text-xs font-normal text-base-content/40">({$projects.length})</span></h3>
        {#if $projects.length === 0}
          <p class="text-sm text-base-content/40 py-4 text-center">暂无项目，请创建一个新项目开始。</p>
        {:else}
          <div class="space-y-1.5">
            {#each $projects as p}
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <!-- svelte-ignore a11y-no-static-element-interactions -->
              <div
                class="flex items-center gap-3 bg-base-300 rounded-lg p-3 cursor-pointer hover:bg-base-300/80 transition-colors group"
                class:ring-1={$currentProject === p.name}
                class:ring-primary={$currentProject === p.name}
                on:click={() => selectProject(p.name)}
              >
                <div class="w-9 h-9 rounded-lg bg-primary/20 text-primary flex items-center justify-center text-sm font-bold shrink-0">
                  {(p.name || '?')[0]}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-medium truncate">{p.name}</div>
                  <div class="text-xs text-base-content/40 truncate">
                    {#if p.title}
                      《{p.title}》
                      {#if p.phase}
                        · {phaseLabels[p.phase] || p.phase}
                      {/if}
                    {:else}
                      空项目
                    {/if}
                  </div>
                </div>
                {#if $currentProject === p.name}
                  <span class="badge badge-primary badge-xs">当前</span>
                {:else}
                  <button
                    class="btn btn-ghost btn-xs text-error opacity-0 group-hover:opacity-100 transition-opacity"
                    on:click|stopPropagation={() => deleteProject(p.name)}
                    disabled={$taskRunning}
                  >
                    删除
                  </button>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
