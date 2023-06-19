#ifndef TTSPIPERLIB_H_
#define TTSPIPERLIB_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

extern void terminate();

extern void initialize();

extern void textToAudio(void* voice, char *text, int cbId);

extern void* loadVoice(char *modelPath, char *modelConfigPath, int64_t *speakerId);

#ifdef __cplusplus
}
#endif

#endif // TTSPIPERLIB_H_