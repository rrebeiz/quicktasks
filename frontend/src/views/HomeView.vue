<template>
  <div class="home">
    <FilterNav @filter="currentFilter = $event" :current="currentFilter" />
    <div v-if="tasks.length">
      <div v-for="task in filteredTasks" :key="task.id">
        <SingleTask
          :task="task"
          @delete="handleDelete"
          @complete="handleComplete"
        />
      </div>
    </div>
    <div v-else>
      <p>No tasks yet, start adding!</p>
    </div>
  </div>
</template>

<script>
import SingleTask from '@/components/SingleTask.vue'
import FilterNav from '@/components/FilterNav.vue'

export default {
  name: 'HomeView',
  data() {
    return {
      tasks: [],
      currentFilter: 'all',
    }
  },
  components: { FilterNav, SingleTask },
  mounted() {
    fetch('http://localhost:4000/v1/tasks')
      .then((res) => res.json())
      .then((data) => {
        if (data.tasks.length) {
          this.tasks = data.tasks
        }
      })
      .catch((err) => console.log(err.message))
  },
  methods: {
    handleDelete(id) {
      this.tasks = this.tasks.filter((task) => task.id !== id)
    },
    handleComplete(id) {
      let t = this.tasks.find((task) => task.id === id)
      t.complete = !t.complete
    },
  },
  computed: {
    filteredTasks() {
      if (this.currentFilter === 'completed') {
        return this.tasks.filter((task) => task.complete)
      }
      if (this.currentFilter === 'ongoing') {
        return this.tasks.filter((task) => !task.complete)
      }
      return this.tasks
    },
  },
}
</script>
