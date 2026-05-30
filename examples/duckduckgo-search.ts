import { searchDuckDuckGo } from '../src/_libs/search/duckduckgo.ts'

const results = await searchDuckDuckGo('cool things')
console.log(JSON.stringify(results, null, 2))
process.exit()
