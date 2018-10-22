module.exports = {
  'root': true,
  'extends': [
    'plugin:vue/recommended',
    'plugin:prettier/recommended',
    'eslint:recommended'
  ],
  'rules': {
    'prettier/prettier': [
      'error',
      {
        'singleQuote': true,
        'semi': false
      }
    ]
  }
}
