import { type Readable, readable } from "svelte/store"

import {
  GetActiveBackups,
  GetActiveMonitors,
  GetActiveRules,
  GetConfig,
  GetErrors,
  GetEvents,
  GetRules,
} from "../wailsjs/go/main/App"
import { EventsOff, EventsOn } from "../wailsjs/runtime"

import type { main } from "../wailsjs/go/models"

type RulePlatform = {
  executable: string
  savegames: string[]
}

export type ActiveRule = {
  ruleFilename: string
  issues: string
  name: string
  platform: RulePlatform
}

export type BackupMetadata = {
  filename: string
  source: string
  backupTime: string
  lastModified: string
}

export const eventStore: Readable<string[]> = readable([], function start(set) {
  async function getEvents() {
    set(await GetEvents())
  }

  getEvents().then(() => {})
  EventsOn("eventsUpdated", function (data) {
    set(data)
  })

  return () => {
    EventsOff("eventsUpdated")
  }
})

export const errorStore: Readable<string[]> = readable([], function start(set) {
  async function getData() {
    set(await GetErrors())
  }

  getData().then(() => {})
  EventsOn("errorsUpdated", function (data) {
    set(data)
  })

  return () => {
    EventsOff("errorsUpdated")
  }
})

export const configStore: Readable<main.Config> = readable(undefined, function start(set) {
  async function getData() {
    set(await GetConfig())
  }

  getData().then(() => {})
  EventsOn("configUpdated", function (data) {
    set(data)
  })

  return () => {
    EventsOff("configUpdated")
  }
})

export const monistorStore: Readable<main.Monitor[]> = readable(undefined, function start(set) {
  async function getData() {
    set(await GetActiveMonitors())
  }

  getData().then(() => {})
  EventsOn("monitorsUpdated", function (data) {
    set(data)
  })

  return () => {
    EventsOff("monitorsUpdated")
  }
})

export const backupStore: Readable<{ [key: string]: BackupMetadata[] }> = readable(
  {},
  function start(set) {
    async function getData() {
      set(await GetActiveBackups())
    }

    getData().then(() => {})
    EventsOn("backupsUpdated", function (data) {
      set(data)
    })

    return () => {
      EventsOff("backupsUpdated")
    }
  }
)

export const activeRuleStore: Readable<{ [key: string]: ActiveRule }> = readable(
  {},
  function start(set) {
    async function getData() {
      set(await GetActiveRules())
    }

    getData().then(() => {})
    EventsOn("activeRulesUpdated", function (data) {
      set(data)
    })

    return () => {
      EventsOff("activeRulesUpdated")
    }
  }
)

export const ruleStore: Readable<{ [key: string]: ActiveRule }> = readable({}, function start(set) {
  async function getData() {
    set(await GetRules())
  }

  getData().then(() => {})
  EventsOn("rulesUpdated", function (data) {
    set(data)
  })

  return () => {
    EventsOff("rulesUpdated")
  }
})
