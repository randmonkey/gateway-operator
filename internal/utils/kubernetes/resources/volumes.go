package resources

import (
	corev1 "k8s.io/api/core/v1"
)

// GetPodVolumeByName gets the pointer of volume with given name.
// if the volume with given name does not exist in the pod, it returns `nil`.
func GetPodVolumeByName(podSpec *corev1.PodSpec, name string) *corev1.Volume {
	for i, volume := range podSpec.Volumes {
		if volume.Name == name {
			return &podSpec.Volumes[i]
		}
	}
	return nil
}
