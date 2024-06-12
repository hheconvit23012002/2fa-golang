<script>

import axiosInstance from "@/utils/request.js";
export default {
  name: "Login",
  data() {
    return {
      phone_number: "",
      password: "",
      submitting:false,
    }
  },
  methods : {
    login()
    {
      this.submitting = true
      axiosInstance.post("/login",{
        phone_number : this.phone_number,
        password : this.password
      }).then((res ) => {
        this.submitting = false
        console.log(res.data.number2FA)
        localStorage.setItem('phoneNumber', this.phone_number);
        localStorage.setItem('auth2FA', res.data.number2FA)
        this.$notify({ type: "success",title: "Login success", text: "thanh cong", duration: 1000 });
        this.$router.push('/2fa');
      }).catch((err) => {
        this.submitting = false
        this.$notify({ type: "error",title: "Login error", text: "that bai" });
        console.log(err)
      })
    }
  }
}
</script>

<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-6">
        <h2 class="mb-4">Đăng Nhập</h2>
        <form method="POST">
          <div class="form-group">
            <label for="exampleInputEmail1">Phone number</label>
            <input type="text" class="form-control" v-model="phone_number" id="exampleInputEmail1" aria-describedby="emailHelp" placeholder="Enter phone" required>
          </div>
          <div class="form-group">
            <label for="exampleInputPassword1">Password</label>
            <input type="password" class="form-control" v-model="password" id="exampleInputPassword1" placeholder="Password" required>
          </div>
          <button type="submit" :disabled="submitting" class="btn btn-primary mt-3" @click.prevent="login">Submit</button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>