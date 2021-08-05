const dict = {
  validations: {
    fields: {
      Title: 'Title',
      Body: 'Body'
    },
    format: (field, message) => `${field}: ${message}`,
    string_eq: p => `length must be equal to ${p}`,
    string_gt: p => `length must be greater than ${p}`,
    string_gte: p => `length must be greater than or equal to ${p}`,
    string_lt: p => `length must be less than ${p}`,
    string_lte: p => `length must be less than or equal to ${p}`,
    string_ne: p => `length must not be ${p}`,
    string_uniqueness: p => `already exists`,
    string_required: p => `must not be blank`
  }
}

function $t (key, fallback = undefined) {
  let items = key.split('.')
  let val = dict
  for (let i = 0; i < items.length; i++) {
    val = val[items[i]]
    if (!val) break
  }
  if (val) return val
  console.debug('MISSING', key)
  if (fallback === undefined) return key
  return fallback
}

export default {
  methods: {
    processErrors (err) {
      let refs = {}
      for (let key in this.$refs) {
        refs[key] = this.$refs[key]
        if (refs[key] && refs[key].$refs) {
          for (let k in refs[key].$refs) {
            refs[k] = refs[key].$refs[k]
          }
        }
      }
      for (let key in refs) {
        let elem = refs[key]
        if (!(elem instanceof Array)) elem = [ elem ]
        elem.forEach((e) => {
          if (e && e.classList) e.classList.remove('is-invalid')
        })
      }
      if (!err || !err.response || !err.response.data) return false
      let errors = err.response.data.Errors
      if (!errors || !errors.length) return false
      let format = $t('validations.format')
      let msgs = []
      errors.forEach(error => {
        let field = $t(`validations.fields.${error.FullName}`, null)
        if (!field) field = $t(`validations.fields.${error.Name}`, error.Name)
        let msg = $t(`validations.${error.Kind}_${error.Type}_${error.Name}`, null)
        if (!msg) msg = $t(`validations.${error.Kind}_${error.Type}`, null)
        if (!msg) msg = $t(`validations.${error.Type}`, error.Type)
        if (typeof(msg) === 'function') {
          msg = msg(error.Param)
        }
        let message = format(field, msg)
        msgs.push(message)
        let key = error.Name
        let elem = refs[key] || refs['input-' + key]
        if (elem) {
          if (elem instanceof Array) elem = elem[0]
          let div = elem.nextElementSibling
          if (!div || !div.classList || !div.classList.contains('invalid-feedback')) {
            let newDiv = document.createElement('div')
            newDiv.classList.add('invalid-feedback')
            elem.parentNode.insertBefore(newDiv, div)
            div = newDiv
          }
          div.innerText = message
          elem.classList.add('is-invalid')
        }
      })
      let firstErrorElem = document.querySelector('.is-invalid')
      if (firstErrorElem) firstErrorElem.focus()
      return true
    }
  }
}
