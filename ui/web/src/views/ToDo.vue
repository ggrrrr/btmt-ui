<template>
  <v-container no-gutters>
    <v-row no-gutters>
      <v-col no-gutters cols="4" sm="4" md="4">
        list
      </v-col>
      <v-col no-gutters>
        <ol-map @pointermove="hoverFeature" :loadTilesWhileAnimating="true" :loadTilesWhileInteracting="true"
          ref="mapRef" style="height: 400px">
          <ol-view ref="view" :center="center" :rotation="rotation" :zoom="zoom" :projection="projection" />
          <ol-tile-layer>
            <ol-source-osm />
          </ol-tile-layer>
          <ol-vector-layer class-name="feature-layer">
            <ol-source-vector url="http://localhost:3000/gara-dolene-ravno-bore.json" :format="geoJson">
              <ol-style>
                <ol-style-stroke color="red" :width="vectorWidth"></ol-style-stroke>
              </ol-style>
            </ol-source-vector>
          </ol-vector-layer>
        </ol-map>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        {{ center }}
        {{ geoJson }}

      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>

import { ref, inject } from 'vue';
// import Map from 'ol/Map';

//[ 25.26545469532615, 41.649423540053725 ]
const vectorWidth = ref(4)
const center = ref([25.26545469532615, 42.649423540053725]);
const projection = ref('EPSG:4326');//EPSG:7796. EPSG:4326/ EPSG:3857
const zoom = ref(8);
const rotation = ref(0);
const format = inject('ol-format');
const geoJson = new format.GeoJSON();
const mapRef = ref(null)

function layerFilter(layerCandidate) {
  console.log('layerFilter', layerCandidate)
  return layerCandidate.getClassName().includes("feature-layer");
}

function hoverFeature(event) {
  const map = mapRef.value?.map
  if (!map) {
    return;
  }

  console.log('hoverFeature.event.pixel', event.pixel)
  // const mapRef = this.$refs['mapRef']
  // console.log('mapRef', map)
  const features = map.getFeaturesAtPixel(event.pixel, {
    hitTolerance: 10,
    layerFilter,
  });

  // console.log('hoverFeature.features', features)
  if (features.length > 0) {
    console.log('hoverFeature.features.0.ol_uid', features[0].ol_uid)
    console.log('hoverFeature.features.0', features[0])
    vectorWidth.value = 8
  }
  //
  // highlightedFeatures.value = features[0] ? [features[0]] : [];
}


</script>
