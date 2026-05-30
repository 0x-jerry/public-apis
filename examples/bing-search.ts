import { searchBing } from '../src/_libs/search/bing.ts'

const results = await searchBing('cool things')
console.log(JSON.stringify(results, null, 2))
process.exit()