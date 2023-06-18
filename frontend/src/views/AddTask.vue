<template>
  <form @submit.prevent="handleSubmit">
    <label for="title">Title:</label>
    <input type="text" v-model="title" required />
    <label for="description">Description:</label>
    <textarea v-model="description" required></textarea>
    <button>Add Task</button>
  </form>
</template>
<script>
export default {
  data() {
    return {
      title: '',
      description: '',
      uri: 'http://localhost:4000/v1/tasks',
    }
  },
  methods: {
    handleSubmit() {
      let task = {
        title: this.title,
        description: this.description,
      }
      fetch(this.uri, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(task),
      })
        .then(() => this.$router.push('/'))
        .catch((err) => console.log(err))
    },
  },
}
</script>

<style>
form {
  background: white;
  padding: 20px;
  border-radius: 10px;
}
label {
  display: block;
  color: #bbb;
  text-transform: uppercase;
  font-size: 14px;
  font-weight: bold;
  letter-spacing: 1px;
  margin: 20px 0 10px 0;
}
input {
  padding: 10px;
  border: 0;
  border-bottom: 1px solid #ddd;
  width: 100%;
  box-sizing: border-box;
}
textarea {
  border: 1px solid #ddd;
  padding: 10px;
  width: 100%;
  box-sizing: border-box;
  height: 100px;
}
form button {
  display: block;
  margin: 20px auto 0;
  background: #00ce89;
  color: white;
  padding: 10px;
  border: 0;
  border-radius: 6px;
  font-size: 16px;
}
</style>
