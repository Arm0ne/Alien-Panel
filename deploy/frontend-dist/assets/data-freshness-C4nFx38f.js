import{$i as e,Bi as t,Ft as n,Gi as r,Jn as i,Ki as a,Li as o,Pn as s,Qr as c,Ri as l,Sa as u,Ui as d,Vi as f,Wi as p,Zi as m,_a as h,ci as g,cr as _,di as v,fi as y,gr as b,ha as x,jn as S,la as C,li as w,mn as T,na as E,pi as D,qn as O,ri as k,ta as A,tr as j,ui as M,yn as N,zi as P,zr as F}from"./router-CvH19P3g.js";import{d as I}from"./index-BTTnI9mu.js";function L(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:i,infoColor:a,successColor:o,warningColor:s,errorColor:l,baseColor:u,borderColor:d,opacityDisabled:f,tagColor:p,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,borderRadiusSmall:_,fontSizeMini:v,fontSizeTiny:y,fontSizeSmall:b,fontSizeMedium:x,heightMini:S,heightTiny:C,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:A}=e;return Object.assign(Object.assign({},I),{closeBorderRadius:_,heightTiny:S,heightSmall:C,heightMedium:w,heightLarge:T,borderRadius:_,opacityDisabled:f,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,fontSizeLarge:x,fontWeightStrong:A,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:u,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:i,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${d}`,textColor:t,color:p,colorBordered:`rgb(250, 250, 252)`,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${c(i,{alpha:.3})}`,textColorPrimary:i,colorPrimary:c(i,{alpha:.12}),colorBorderedPrimary:c(i,{alpha:.1}),closeIconColorPrimary:i,closeIconColorHoverPrimary:i,closeIconColorPressedPrimary:i,closeColorHoverPrimary:c(i,{alpha:.12}),closeColorPressedPrimary:c(i,{alpha:.18}),borderInfo:`1px solid ${c(a,{alpha:.3})}`,textColorInfo:a,colorInfo:c(a,{alpha:.12}),colorBorderedInfo:c(a,{alpha:.1}),closeIconColorInfo:a,closeIconColorHoverInfo:a,closeIconColorPressedInfo:a,closeColorHoverInfo:c(a,{alpha:.12}),closeColorPressedInfo:c(a,{alpha:.18}),borderSuccess:`1px solid ${c(o,{alpha:.3})}`,textColorSuccess:o,colorSuccess:c(o,{alpha:.12}),colorBorderedSuccess:c(o,{alpha:.1}),closeIconColorSuccess:o,closeIconColorHoverSuccess:o,closeIconColorPressedSuccess:o,closeColorHoverSuccess:c(o,{alpha:.12}),closeColorPressedSuccess:c(o,{alpha:.18}),borderWarning:`1px solid ${c(s,{alpha:.35})}`,textColorWarning:s,colorWarning:c(s,{alpha:.15}),colorBorderedWarning:c(s,{alpha:.12}),closeIconColorWarning:s,closeIconColorHoverWarning:s,closeIconColorPressedWarning:s,closeColorHoverWarning:c(s,{alpha:.12}),closeColorPressedWarning:c(s,{alpha:.18}),borderError:`1px solid ${c(l,{alpha:.23})}`,textColorError:l,colorError:c(l,{alpha:.1}),colorBorderedError:c(l,{alpha:.08}),closeIconColorError:l,closeIconColorHoverError:l,closeIconColorPressedError:l,closeColorHoverError:c(l,{alpha:.12}),closeColorPressedError:c(l,{alpha:.18})})}var R={name:`Tag`,common:T,self:L},z={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},B=w(`tag`,`
 --n-close-margin: var(--n-close-margin-top) var(--n-close-margin-right) var(--n-close-margin-bottom) var(--n-close-margin-left);
 white-space: nowrap;
 position: relative;
 box-sizing: border-box;
 cursor: default;
 display: inline-flex;
 align-items: center;
 flex-wrap: nowrap;
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 line-height: 1;
 height: var(--n-height);
 font-size: var(--n-font-size);
`,[v(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),M(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),M(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),M(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),M(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),v(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[M(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),M(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),v(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),v(`icon, avatar`,[v(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),v(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),v(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[y(`disabled`,[g(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[y(`checked`,`color: var(--n-text-color-hover-checkable);`)]),g(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[y(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),v(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[y(`disabled`,[g(`&:hover`,`background-color: var(--n-color-checked-hover);`),g(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),V=Object.assign(Object.assign(Object.assign({},S.props),z),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),H=F(`n-tag`),U=r({name:`Tag`,props:V,slots:Object,setup(e){let t=x(null),{mergedBorderedRef:n,mergedClsPrefixRef:r,inlineThemeDisabled:a,mergedRtlRef:c,mergedComponentPropsRef:l}=i(e),u=o(()=>e.size||l?.value?.Tag?.size||`medium`),d=S(`Tag`,`-tag`,B,R,e,r);E(H,{roundRef:h(e,`round`)});function f(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function p(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&_(n,t)}}let m={setTextContent(e){let{value:n}=t;n&&(n.textContent=e)}},g=s(`Tag`,c,r),v=o(()=>{let{type:t,color:{color:r,textColor:i}={}}=e,a=u.value,{common:{cubicBezierEaseInOut:o},self:{padding:s,closeMargin:c,borderRadius:l,opacityDisabled:f,textColorCheckable:p,textColorHoverCheckable:m,textColorPressedCheckable:h,textColorChecked:g,colorCheckable:_,colorHoverCheckable:v,colorPressedCheckable:y,colorChecked:b,colorCheckedHover:x,colorCheckedPressed:S,closeBorderRadius:C,fontWeightStrong:w,[D(`colorBordered`,t)]:T,[D(`closeSize`,a)]:E,[D(`closeIconSize`,a)]:O,[D(`fontSize`,a)]:A,[D(`height`,a)]:j,[D(`color`,t)]:M,[D(`textColor`,t)]:N,[D(`border`,t)]:P,[D(`closeIconColor`,t)]:F,[D(`closeIconColorHover`,t)]:I,[D(`closeIconColorPressed`,t)]:L,[D(`closeColorHover`,t)]:R,[D(`closeColorPressed`,t)]:z}}=d.value,B=k(c);return{"--n-font-weight-strong":w,"--n-avatar-size-override":`calc(${j} - 8px)`,"--n-bezier":o,"--n-border-radius":l,"--n-border":P,"--n-close-icon-size":O,"--n-close-color-pressed":z,"--n-close-color-hover":R,"--n-close-border-radius":C,"--n-close-icon-color":F,"--n-close-icon-color-hover":I,"--n-close-icon-color-pressed":L,"--n-close-icon-color-disabled":F,"--n-close-margin-top":B.top,"--n-close-margin-right":B.right,"--n-close-margin-bottom":B.bottom,"--n-close-margin-left":B.left,"--n-close-size":E,"--n-color":r||(n.value?T:M),"--n-color-checkable":_,"--n-color-checked":b,"--n-color-checked-hover":x,"--n-color-checked-pressed":S,"--n-color-hover-checkable":v,"--n-color-pressed-checkable":y,"--n-font-size":A,"--n-height":j,"--n-opacity-disabled":f,"--n-padding":s,"--n-text-color":i||N,"--n-text-color-checkable":p,"--n-text-color-checked":g,"--n-text-color-hover-checkable":m,"--n-text-color-pressed-checkable":h}}),y=a?O(`tag`,o(()=>{let t=``,{type:r,color:{color:i,textColor:a}={}}=e;return t+=r[0],t+=u.value[0],i&&(t+=`a${b(i)}`),a&&(t+=`b${b(a)}`),n.value&&(t+=`c`),t}),v,e):void 0;return Object.assign(Object.assign({},m),{rtlEnabled:g,mergedClsPrefix:r,contentRef:t,mergedBordered:n,handleClick:f,handleCloseClick:p,cssVars:a?void 0:v,themeClass:y?.themeClass,onRender:y?.onRender})},render(){var e;let{mergedClsPrefix:t,rtlEnabled:n,closable:r,color:{borderColor:i}={},round:o,onRender:s,$slots:c}=this;s?.();let l=j(c.avatar,e=>e&&a(`div`,{class:`${t}-tag__avatar`},e)),u=j(c.icon,e=>e&&a(`div`,{class:`${t}-tag__icon`},e));return a(`div`,{class:[`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:o,[`${t}-tag--avatar`]:l,[`${t}-tag--icon`]:u,[`${t}-tag--closable`]:r}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},u||l,a(`span`,{class:`${t}-tag__content`,ref:`contentRef`},(e=this.$slots).default?.call(e)),!this.checkable&&r?a(N,{clsPrefix:t,class:`${t}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:o,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?a(`div`,{class:`${t}-tag__border`,style:{borderColor:i}}):null)}}),W={key:0,class:`ml-4px`},G={key:0,class:`mt-2px text-12px opacity-80`},K=r({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(r){let i=r,a=x(Date.now()),s,c=o(()=>{if(!i.dataAt)return null;let e=new Date(i.dataAt).getTime();return Number.isFinite(e)?e:null}),h=o(()=>{if(c.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=a.value-c.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(i.delayedAfterMinutes,1)*60*1e3,n=Math.max(i.expiredAfterMinutes,i.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),g=o(()=>c.value===null?``:new Date(c.value).toLocaleString(`zh-CN`,{hour12:!1}));return e(()=>{i.compact||(s=setInterval(()=>{a.value=Date.now()},6e4))}),m(()=>{s&&clearInterval(s)}),(e,i)=>{let a=U,o=n;return A(),P(o,null,{trigger:C(()=>[p(a,{size:`small`,type:h.value.type},{default:C(()=>[d(u(h.value.label)+` `,1),!r.compact&&g.value?(A(),f(`span`,W,`· `+u(g.value),1)):t(``,!0)]),_:1},8,[`type`])]),default:C(()=>[l(`div`,null,u(h.value.detail),1),g.value?(A(),f(`div`,G,`数据时间：`+u(g.value),1)):t(``,!0)]),_:1})}}});export{U as n,K as t};