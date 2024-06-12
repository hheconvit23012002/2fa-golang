<script>
import axiosInstance from "@/utils/request.js";

export default {
  name: "AuthTwoFA",
  data() {
    return {
      seconds_remaining : 40,
      number : null,
      phone_number : null,
      timer : null
    }
  },
  mounted() {
    this.getLocalStorage();
    this.callToCustomer();
    this.startTimer();
  },
  methods : {
    callToCustomer() {
      axiosInstance.post("/callToCustomer",{
        phone_number : this.phone_number
      }).then((res) => {
        this.$notify({ type: "success",title: "Call success", text: "thanh cong", duration: 1000 });
      }).catch((error) => {
        this.$notify({ type: "error",title: "Call error", text: "that bai" });
        console.log(error)
      })
    },
    getLocalStorage() {
      this.number = localStorage.getItem('auth2FA');
      this.phone_number = localStorage.getItem('phoneNumber');
    },
    check2FA(){
      axiosInstance.post("/check2FA",{
        phone_number: this.phone_number
      }).then((res ) => {
        console.log(res)
        localStorage.removeItem('phoneNumber');
        localStorage.removeItem('auth2FA')
        clearInterval(this.timer)
        this.$notify({ type: "success",title: "Login success", text: "thanh cong", duration: 1000 });
        this.$router.push('/loginSuccess');
      }).catch((err) => {
        console.log(err)
      })
    },
    startTimer() {
       this.timer = setInterval(() => {
        if(this.seconds_remaining < 0){
          clearInterval(this.timer)
          this.$router.push('/');
        }
        this.check2FA();
        this.seconds_remaining-=1
      }, 1000);
    },
  }
}
</script>

<template>
  <div class="container mt-5">
    <div class="card">
      <div class="card-body text-center">
        <h5 class="card-title">Số điện thoại của bạn là : {{ this.phone_number}}</h5>
        <p class="card-text">Số bạn cần nhập xác thực 2 lớp là : {{ this.number}}</p>
        <div class="account-number">
          <span class="badge bg-primary p-3">{{ this.seconds_remaining}}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>