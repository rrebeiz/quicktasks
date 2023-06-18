<template>
  <h1>Edit task</h1>
  <form @submit.prevent="handleSubmit">
    <label for="title">Title:</label>
    <input type="text" v-model="title" required />
    <label for="description">Description:</label>
    <textarea v-model="description" required></textarea>
    <button>Update Task</button>
  </form>
</template>
<script>
export default {
  name: 'EditTask',
  props: ['id'],
  data() {
    return {
      title: '',
      description: '',
      uri: 'http://localhost:4000/v1/tasks/' + this.id,
    }
  },
  mounted() {
    fetch(this.uri)
      .then((res) => res.json())
      .then((data) => {
        this.title = data.task.title
        this.description = data.task.description
      })
      .catch((err) => console.log(err))
  },
  methods: {
    handleSubmit() {
      fetch(this.uri, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: this.title,
          description: this.description,
        }),
      })
        .then(() => {
          this.$router.push('/')
        })
        .catch((err) => console.log(err))
    },
  },
}
</script>

<style></style>
