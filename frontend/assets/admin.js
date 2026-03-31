(function () {
  const api = window.blogApi;
  const auth = window.blogAuth;

  const state = {
    editingSlug: null,
    categories: [],
    tags: [],
    page: 1,
    pageSize: 20,
  };

  const $ = (id) => document.getElementById(id);

  function showMsg(type, text) {
    $('error').textContent = type === 'error' ? text : '';
    $('success').textContent = type === 'success' ? text : '';
  }

  function setLoginState() {
    const loggedIn = !!auth.getToken();
    $('login-panel').classList.toggle('hidden', loggedIn);
    $('admin-panel').classList.toggle('hidden', !loggedIn);
    $('logout-btn').classList.toggle('hidden', !loggedIn);
  }

  function switchTab(tab) {
    document.querySelectorAll('[data-tab]').forEach((el) => el.classList.add('hidden'));
    document.querySelectorAll('.tab-btn').forEach((el) => el.classList.remove('active'));
    $(`tab-${tab}`).classList.add('active');
    $(`panel-${tab}`).classList.remove('hidden');
  }

  async function doLogin() {
    try {
      const username = $('username').value.trim();
      const password = $('password').value;
      const out = await api.login(username, password);
      auth.setToken(out.token);
      setLoginState();
      await loadBootData();
      showMsg('success', '登录成功');
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  function fillTaxonomySelects() {
    const catSelect = $('post-category');
    catSelect.innerHTML = '<option value="">无分类</option>';
    state.categories.forEach((c) => {
      const op = document.createElement('option'); op.value = c.id; op.textContent = c.name; catSelect.appendChild(op);
    });

    const tagsWrap = $('post-tags');
    tagsWrap.innerHTML = '';
    state.tags.forEach((t) => {
      const label = document.createElement('label');
      label.style.marginRight = '10px';
      label.innerHTML = `<input type="checkbox" value="${t.id}"> ${t.name}`;
      tagsWrap.appendChild(label);
    });
  }

  function selectedTagIDs() {
    return [...$('post-tags').querySelectorAll('input[type=checkbox]:checked')].map((el) => Number(el.value));
  }

  async function loadBootData() {
    const [cats, tags] = await Promise.all([api.listCategories(), api.listTags()]);
    state.categories = cats.data || [];
    state.tags = tags.data || [];
    fillTaxonomySelects();
    await Promise.all([renderCategories(), renderTags(), loadPosts()]);
  }

  async function loadPosts() {
    try {
      const params = {
        page: String(state.page),
        page_size: String(state.pageSize),
        q: $('post-q').value,
        published: $('post-published').value,
      };
      const res = await api.listAdminPosts(params);
      $('post-total').textContent = `共 ${res.total || 0} 篇`;
      const body = $('post-table-body');
      body.innerHTML = '';
      for (const p of (res.data || [])) {
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${p.title}</td>
          <td>${p.published ? '已发布' : '草稿'}</td>
          <td>${p.category?.name || '-'}</td>
          <td>${new Date(p.created_at).toLocaleDateString()}</td>
          <td>
            <button data-edit="${p.slug}" class="secondary">编辑</button>
            <button data-del="${p.slug}" class="danger">删除</button>
          </td>`;
        body.appendChild(tr);
      }
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function renderCategories() {
    const list = $('categories-list');
    list.innerHTML = '';
    for (const c of state.categories) {
      const item = document.createElement('div');
      item.className = 'card';
      item.innerHTML = `<strong>${c.name}</strong> <span class='small'>/${c.slug}</span>
      <button class='danger' data-del-cat='${c.id}' style='float:right'>删除</button>`;
      list.appendChild(item);
    }
  }

  async function renderTags() {
    const list = $('tags-list');
    list.innerHTML = '';
    for (const t of state.tags) {
      const item = document.createElement('div');
      item.className = 'card';
      item.innerHTML = `<strong>${t.name}</strong> <span class='small'>/${t.slug}</span>
      <button class='danger' data-del-tag='${t.id}' style='float:right'>删除</button>`;
      list.appendChild(item);
    }
  }

  function resetEditor() {
    state.editingSlug = null;
    $('editor-title').textContent = '新建文章';
    $('post-title').value = '';
    $('post-cover').value = '';
    $('post-content').value = '';
    $('post-category').value = '';
    $('post-published').checked = false;
    $('upload-file').value = '';
    $('post-tags').querySelectorAll('input[type=checkbox]').forEach((cb) => cb.checked = false);
  }

  async function editPost(slug) {
    try {
      const { data: p } = await api.getAnyPost(slug);
      state.editingSlug = slug;
      $('editor-title').textContent = `编辑：${p.title}`;
      $('post-title').value = p.title || '';
      $('post-cover').value = p.cover_url || '';
      $('post-content').value = p.content || '';
      $('post-category').value = p.category_id || '';
      $('post-published').checked = !!p.published;
      const selected = new Set((p.tags || []).map(t => t.id));
      $('post-tags').querySelectorAll('input[type=checkbox]').forEach((cb) => {
        cb.checked = selected.has(Number(cb.value));
      });
      switchTab('editor');
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function savePost() {
    const payload = {
      title: $('post-title').value.trim(),
      content: $('post-content').value,
      cover_url: $('post-cover').value.trim(),
      category_id: $('post-category').value ? Number($('post-category').value) : null,
      tag_ids: selectedTagIDs(),
      published: $('post-published').checked,
    };

    try {
      if (state.editingSlug) {
        await api.updatePost(state.editingSlug, payload);
        showMsg('success', '文章已更新');
      } else {
        await api.createPost(payload);
        showMsg('success', '文章已创建');
      }
      resetEditor();
      switchTab('posts');
      await loadPosts();
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function removePost(slug) {
    if (!confirm(`确认删除 ${slug} ?`)) return;
    try {
      await api.deletePost(slug);
      showMsg('success', '文章已删除');
      await loadPosts();
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function uploadImage() {
    const file = $('upload-file').files[0];
    if (!file) return;
    try {
      const res = await api.uploadImage(file);
      $('post-cover').value = res.url;
      showMsg('success', '图片上传成功');
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function createCategory() {
    const name = $('new-category').value.trim();
    if (!name) return;
    try {
      await api.createCategory(name);
      $('new-category').value = '';
      await loadBootData();
      showMsg('success', '分类已创建');
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function createTag() {
    const name = $('new-tag').value.trim();
    if (!name) return;
    try {
      await api.createTag(name);
      $('new-tag').value = '';
      await loadBootData();
      showMsg('success', '标签已创建');
    } catch (e) {
      showMsg('error', e.message);
    }
  }

  async function deleteCategory(id) {
    if (!confirm('确认删除该分类？')) return;
    await api.deleteCategory(id);
    await loadBootData();
  }

  async function deleteTag(id) {
    if (!confirm('确认删除该标签？')) return;
    await api.deleteTag(id);
    await loadBootData();
  }

  function bindEvents() {
    $('login-btn').addEventListener('click', doLogin);
    $('logout-btn').addEventListener('click', () => {
      auth.setToken('');
      setLoginState();
      showMsg('success', '已退出登录');
    });

    $('tab-posts').addEventListener('click', () => switchTab('posts'));
    $('tab-editor').addEventListener('click', () => switchTab('editor'));
    $('tab-categories').addEventListener('click', () => switchTab('categories'));
    $('tab-tags').addEventListener('click', () => switchTab('tags'));

    $('search-posts').addEventListener('click', () => { state.page = 1; loadPosts(); });
    $('prev-page').addEventListener('click', () => { state.page = Math.max(1, state.page - 1); $('page-now').textContent = state.page; loadPosts(); });
    $('next-page').addEventListener('click', () => { state.page += 1; $('page-now').textContent = state.page; loadPosts(); });

    $('save-post').addEventListener('click', savePost);
    $('new-post').addEventListener('click', () => { resetEditor(); switchTab('editor'); });
    $('upload-btn').addEventListener('click', uploadImage);

    $('create-category').addEventListener('click', createCategory);
    $('create-tag').addEventListener('click', createTag);

    $('post-table-body').addEventListener('click', (e) => {
      const edit = e.target.getAttribute('data-edit');
      const del = e.target.getAttribute('data-del');
      if (edit) editPost(edit);
      if (del) removePost(del);
    });

    $('categories-list').addEventListener('click', (e) => {
      const id = e.target.getAttribute('data-del-cat');
      if (id) deleteCategory(id).catch((err) => showMsg('error', err.message));
    });
    $('tags-list').addEventListener('click', (e) => {
      const id = e.target.getAttribute('data-del-tag');
      if (id) deleteTag(id).catch((err) => showMsg('error', err.message));
    });
  }

  async function init() {
    bindEvents();
    setLoginState();
    switchTab('posts');
    $('page-now').textContent = state.page;
    if (auth.getToken()) {
      try {
        await loadBootData();
      } catch (e) {
        showMsg('error', e.message);
      }
    }
  }

  init();
})();
