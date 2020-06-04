const getHost = () => {
  console.log(process.env.backend)
    return process.env.backend || window.location.host
}

export default ({},inject) => {
  inject("getHost",getHost);
}
