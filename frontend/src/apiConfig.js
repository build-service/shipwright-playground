let apiUrl
const apiUrls = {
  production: 'x',
  development: 'x'
}

if (window.location.hostname === 'localhost') {
  apiUrl = apiUrls.development
} else {
  apiUrl = apiUrls.production
}

export default apiUrl